package cdp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/visualiser"
)

var (
	testOutFolder = "testOutFolder/"
)

func init() {
	cfg.SetConfigFileName("")
	viper.Set(cfg.DataDir, testOutFolder)
}

type testSz struct {
	name string
	user szenario.User
	exec szenario.JobFunc
}

// Name returns the name
func (s testSz) Name() string {
	return s.name
}

func (s testSz) SetName(string) {}

// GetUser returns the user the szenario runs with
func (s *testSz) User() szenario.User {
	return s.user
}

// SetUser set the user the szenario runs with
func (s *testSz) SetUser(u szenario.User) {
	s.user = u
}

// Execute the szenario
func (s *testSz) Execute(cdp szenario.Engine) (err error) {
	return s.exec(cdp)
}

func NewSzenario(n string, exec szenario.JobFunc) szenario.Szenario {
	return &testSz{name: n, exec: exec}
}

func initEnv(t *testing.T) (*core.Bus, []Option, func()) {
	opts := make([]Option, 0)
	opts = append(opts, Timeout(120*time.Second))

	hcl := hcl.New(hcl.WithName(t.Name()), hcl.WithLevel(hclog.Warn))
	c, close := core.New("som-test", core.HCL(hcl))
	bus := c.Bus()
	visualiser.NewDumper()

	return bus, opts, func() {
		close()
		cleanupOutFolder()
	}
}

func cleanupOutFolder() {
	if err := os.RemoveAll(testOutFolder); err != nil {
		hcl.Warnf("Cannot remove %s: %v", testOutFolder, err)
	}
}

func TestTimeOut(t *testing.T) {
	bus, opts, cleanupEnv := initEnv(t)
	defer cleanupEnv()
	opts = append(opts, Timeout(2*time.Second))

	htmlData := "MyAnswer"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, htmlData)
	}))
	defer srv.Close()
	timeout := false
	cdp, cancel := New(opts...)
	defer cancel()

	sz := NewSzenario("TestTimeOut",
		func(cdp szenario.Engine) error {
			cdp.Step("Open", chromedp.Navigate(srv.URL))
			cdp.Step("Check data", cdp.Body(cdp.Contains(htmlData)))
			cdp.WaitForEver()
			return nil
		},
	)
	bus.Szenario.Handle(func(e *msg.SzenarioEvtMsg) {
		if e.Err() != nil {
			timeout = true
		}
	})

	cdp.Execute(sz)
	bus.Szenario.WaitMsgProcessed()
	if !timeout {
		t.Error("Timeout is not triggred")
	}
	files, err := ioutil.ReadDir(testOutFolder)
	if err != nil {
		t.Errorf("Cannot read output dir %s: %v", testOutFolder, err)
	}
	if len(files) != 2 {
		t.Error("Unexpected number of output files:")
		for _, f := range files {
			t.Errorf("file %v", f.Name())
		}
	}
}

func TestBodyDump(t *testing.T) {
	bus, opts, cleanupEnv := initEnv(t)
	defer cleanupEnv()
	htmlBody := `Title, Some data`
	htmlData := fmt.Sprintf(`<html>
<body>
%s
</body>
</html>`, htmlBody)
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, htmlData)
		called = true
	}))
	defer srv.Close()
	var srvErr error
	cdp, cancel := New(opts...)
	defer cancel()
	sz := NewSzenario("TestBody",
		func(cdp szenario.Engine) error {
			cdp.Step("Open", chromedp.Navigate(srv.URL))
			var content string
			cdp.Step("Strings", cdp.Body(cdp.Strings(&content)))
			if !strings.Contains(content, htmlBody) {
				t.Errorf("Body saved to string: %q should be %q", content, content)
			}
			return nil
		},
	)

	bus.Szenario.Handle(func(e *msg.SzenarioEvtMsg) {
		if e.Err() != nil {
			srvErr = e.Err()
		}
	})

	cdp.Execute(sz)
	bus.Szenario.WaitMsgProcessed()

	if !called {
		t.Error("Server was not called")
	}

	if srvErr != nil {
		t.Errorf("Script should not error: %v", srvErr)
	}
}

func getCounter(couters map[string]any, key string) int {
	v, ok := couters[key]
	if !ok {
		return -1
	}
	if i, ok := v.(int); ok {
		return i
	}
	if f, ok := v.(float64); ok {
		return int(f)
	}

	return -1
}

func TestConsoleLog(t *testing.T) {
	bus, opts, cleanupEnv := initEnv(t)
	defer cleanupEnv()
	htmlData := `<html>
<body>
<script>
console.warn("TestWarning")
console.warn("TestWarning")
console.warn("TestWarning")
console.warn("TestWarning")
console.error("TestError")
console.error("TestError")
console.error("TestError")
console.assert(false, "testAssert")
console.assert(false, "testAssert")
trow("testException")
</script>
<h1 id="h1id">title</h1>
</body>
</html>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, htmlData)
	}))
	defer srv.Close()
	var srvErr error
	cdp, cancel := New(opts...)
	defer cancel()
	sz := NewSzenario("TestBody",
		func(cdp szenario.Engine) error {
			cdp.Step("Open", chromedp.Navigate(srv.URL))
			cdp.Step("Wait", chromedp.WaitVisible("#h1id", chromedp.ByID))
			return nil
		},
	)

	called := false
	bus.Szenario.Handle(func(e *msg.SzenarioEvtMsg) {
		if e.Err() != nil {
			srvErr = e.Err()
		}
		called = true
		ctr := e.Counters
		if getCounter(ctr, "console.assert") != 2 {
			t.Errorf("Should have 2 asserts")
		}
		if getCounter(ctr, "console.error") != 3 {
			t.Errorf("Should have 2 errors")
		}
		if getCounter(ctr, "console.warning") != 4 {
			t.Errorf("Should have 2 warnings")
		}
		if getCounter(ctr, "console.exception") != 1 {
			t.Errorf("Should have 1 exceptions")
		}
	})

	cdp.Execute(sz)
	bus.Szenario.WaitMsgProcessed()

	if !called {
		t.Error("Server was not called")
	}
	if srvErr != nil {
		t.Errorf("Script should not error: %v", srvErr)
	}
}

type bodyTestCase struct {
	name    string
	body    string
	action  chromedp.Action
	wantErr bool
}

func TestBody(t *testing.T) {
	bus, opts, cleanupEnv := initEnv(t)
	defer cleanupEnv()
	cdp, cancel := New(opts...)
	defer cancel()
	tests := []bodyTestCase{
		{
			name:    "Check text",
			body:    "<h1>Title</h1>\n Some data",
			action:  cdp.Body(cdp.Contains("Some data")),
			wantErr: false,
		},
		{
			name:    "Check text - but not found",
			body:    "<h1>Title</h1>\n Some data",
			action:  cdp.Body(cdp.Contains("Not there")),
			wantErr: true,
		},
		{
			name:    "Check title",
			body:    "<h1>Title</h1>\n Some data",
			action:  cdp.Body(cdp.Contains("Title")),
			wantErr: false,
		},
		{
			name:    "Check not to be found",
			body:    "<h1>Title</h1>\n Some data",
			action:  cdp.Body(cdp.NotContains("not to be found")),
			wantErr: false,
		},
		{
			name:    "Check not to be found, but there",
			body:    "<h1>Title</h1>\n Some data",
			action:  cdp.Body(cdp.NotContains("Some data")),
			wantErr: true,
		},
		{
			name:    "Check size",
			body:    "Title, Some data",
			action:  cdp.Body(cdp.Bigger(len("Title, Some data") - 1)),
			wantErr: false,
		},
		{
			name:    "Check size, but too small",
			body:    "Title, Some data",
			action:  cdp.Body(cdp.Bigger(len("Title, Some data") + 1)),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		runTestBody(t, bus, cdp, &tc)
	}

}

func runTestBody(t *testing.T, bus *core.Bus, cdp *Engine, tc *bodyTestCase) {
	core.EnsureOutFolder(testOutFolder)
	defer cleanupOutFolder()
	var srvErr error
	htmlData := fmt.Sprintf(`<html>
<body>
%s
</body>
</html>`, tc.body)
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, htmlData)
		called = true
	}))
	defer srv.Close()
	sz := NewSzenario("TestBody",
		func(cdp szenario.Engine) error {
			cdp.Step("Open", chromedp.Navigate(srv.URL))
			cdp.Step(tc.name, tc.action)
			// if (srvErr != nil) != tc.wantErr {
			// 	t.Errorf("Want err %v got: %v", tc.wantErr, srvErr)
			// }
			return nil
		},
	)

	bus.Szenario.Handle(func(e *msg.SzenarioEvtMsg) {
		if e.Err() != nil {
			srvErr = e.Err()
		}
	})
	cdp.runChan = make(chan szenarionRunWrapper, 100)
	cdp.Execute(sz)
	bus.Szenario.WaitMsgProcessed()
	if !called {
		t.Error("Server was not called")
	}

	if (srvErr != nil) != tc.wantErr {
		t.Errorf("Script %q should error %v got: %v", tc.name, tc.wantErr, srvErr)
	}

	files, err := ioutil.ReadDir(testOutFolder)
	if err != nil {
		t.Errorf("Cannot read output dir %s: %v", testOutFolder, err)
	}
	fileCnt := 0
	if tc.wantErr {
		fileCnt = 2
		// time.Sleep(500 * time.Millisecond)
	}
	if len(files) != fileCnt {
		t.Error("Unexpected number of output files:")
		for _, f := range files {
			t.Errorf("file %v", f.Name())
		}
	}
}

func TestReporter(t *testing.T) {
	bus, opts, cleanupEnv := initEnv(t)
	defer cleanupEnv()
	htmlData := `<html>
<body>
<h1 id="h1id">title</h1>
</body>
</html>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, htmlData)
	}))
	ncDummy := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// go does not like its own self signed certs
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		s := string(b)
		hcl.Warn(s)
	}))
	ncDummy.StartTLS()
	defer srv.Close()
	var srvErr error
	cdp, cancel := New(opts...)
	defer cancel()
	sz := NewSzenario("TestBody",
		func(cdp szenario.Engine) error {
			cdp.Step("Open", chromedp.Navigate(srv.URL))
			cdp.Step("Wait", chromedp.WaitVisible("#h1id", chromedp.ByID))
			return nil
		},
	)
	// bridger.RegisterNetCrunchWebMessage(sz.Name(), sz.GetBackends())

	bus.Szenario.Handle(func(e *msg.SzenarioEvtMsg) {
		if e.Err() != nil {
			srvErr = e.Err()
		}
	})
	cdp.Execute(sz)
	bus.Szenario.WaitMsgProcessed()

	if srvErr != nil {
		t.Errorf("Script should not error: %v", srvErr)
	}
}

func TestRepeat(t *testing.T) {
	bus, opts, cleanupEnv := initEnv(t)
	defer cleanupEnv()
	opts = append(opts, Repeat(1*time.Second), Timeout(3*time.Second))
	htmlData := `<html>
<body>
<h1 id="h1id">title</h1>
</body>
</html>`
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, htmlData)
		calls++
	}))
	cdp, cancel := New(opts...)
	cdp.repeat = 1 * time.Second
	sz := NewSzenario("TestBody",
		func(cdp szenario.Engine) error {
			cdp.Step("Open", chromedp.Navigate(srv.URL))
			cdp.Step("Wait", chromedp.WaitVisible("#h1id", chromedp.ByID))
			return nil
		},
	)
	end := make(chan bool)
	evtCalls := 0
	bus.Szenario.Handle(func(e *msg.SzenarioEvtMsg) {
		if e.Err() == nil {
			evtCalls++
		}
		if evtCalls > 2 {
			end <- true
		}
	})
	timer := time.AfterFunc(50*time.Second, func() {
		end <- true
	})
	go cdp.Execute(sz)
	<-end
	timer.Stop()
	bus.Szenario.WaitMsgProcessed()
	cancel()

	if calls < 3 {
		t.Errorf("did not repeat: calls %v", calls)
	}

}
