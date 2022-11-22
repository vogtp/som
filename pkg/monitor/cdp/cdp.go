package cdp

import (
	"context"
	"errors"
	"fmt"

	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/mime"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater/user"
)

// Option configures the Engine
type Option func(*Engine)

// Engine contains all information to run a chrome szenario
type Engine struct {
	ctx     context.Context
	baseHcl hcl.Logger
	hcl     hcl.Logger
	bus     *core.Bus

	runChan        chan szenarionRunWrapper
	stepBreakPoint chan any

	// muStep protects step related stuff
	muStep   sync.Mutex
	browser  context.Context
	szenario szenario.Szenario
	stepInfo *stepInfo
	consMsg  map[string]int

	// mu protects internals
	mu     sync.Mutex
	evtMsg *msg.SzenarioEvtMsg

	// flags
	headless      bool
	noClose       bool
	timeout       time.Duration
	stepDelay     time.Duration
	timeoutTicker *time.Ticker
	sendReport    bool
}

// New creates a new Engine
func New(opts ...Option) (*Engine, context.CancelFunc) {
	core := core.Get()
	hcl := core.HCL().Named("cdp")
	ctx, cancel := context.WithCancel(context.Background())
	cdp := &Engine{
		ctx:        ctx,
		baseHcl:    hcl,
		hcl:        hcl,
		bus:        core.Bus(),
		runChan:    make(chan szenarionRunWrapper, 100),
		sendReport: true,
		consMsg:    make(map[string]int),
		headless:   !viper.GetBool(cfg.BrowserShow),
		noClose:    viper.GetBool(cfg.BrowserNoClose),
		timeout:    viper.GetDuration(cfg.CheckTimeout),
		stepDelay:  viper.GetDuration(cfg.CheckStepDelay),
	}

	for _, o := range opts {
		o(cdp)
	}

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-signalChanel
		cdp.hcl.Print("Canceled by OS signal")
		cdp.sendReport = false
		cancel()
		os.Exit(1)
	}()

	return cdp, cancel
}

// SetTimeout updates the timeout and restarts the timer
func (cdp *Engine) SetTimeout(d time.Duration) {
	cdp.mu.Lock()
	defer cdp.mu.Unlock()
	cdp.timeout = d
	if cdp.timeoutTicker != nil {
		cdp.hcl.Infof("Updating timeout %v", d)
		cdp.timeoutTicker.Reset(d)
	}
}

type szenarionRunWrapper struct {
	sz        szenario.Szenario
	lastRunOk bool
	retry     int
}

// Execute runs one or more szenarios
func (cdp *Engine) RunUser(username string) error {
	user, err := user.Store.Get(username)
	if err != nil {
		return fmt.Errorf("no such user: %v", username)
	}
	sc := core.Get().SzenaioConfig()
	if sc == nil || sc == szenario.NoConfig {
		return fmt.Errorf("No szenarion config loaded!")
	}
	szs, err := sc.ByUser(user)
	if err != nil {
		return fmt.Errorf("cannot get szenarios for user %v", username)
	}
	if viper.GetBool(cfg.PasswdChange) {
		go cdp.passwordChangeLoop(user)
	}
	cdp.schedule(szs...)
	cdp.loop()
	return nil
}

// Execute runs one or more szenarios
func (cdp *Engine) Execute(szenarios ...szenario.Szenario) {
	cdp.schedule(szenarios...)
	cdp.loop()
}

// schedule one or more szenarios
func (cdp *Engine) schedule(szenarios ...szenario.Szenario) {
	cdp.hcl.Infof("Scheduling %d szenarios.", len(szenarios))
	for _, sz := range szenarios {
		cdp.hcl.Infof("Scheduling szenario %s", sz.Name())
		cdp.runChan <- szenarionRunWrapper{sz: sz}
	}
}

// loop runs the szenarios and never returns
func (cdp *Engine) loop() {
	for srw := range cdp.runChan {
		cdp.hcl = cdp.baseHcl.Named(srw.sz.Name())
		cdp.szenario = srw.sz
		cdp.stepInfo = newStepInfo(&cdp.hcl)
		srw.lastRunOk = cdp.run()
		cdp.hcl = cdp.baseHcl
		go cdp.reschedule(srw)
	}
	cdp.hcl.Info("Finished main szenario loop")
	if cdp.noClose {
		<-cdp.browser.Done()
		return
	}
}

func (cdp *Engine) rescheduleDelay(srw *szenarionRunWrapper) time.Duration {
	delay := srw.sz.RepeatDelay()
	if srw.lastRunOk {
		return delay
	}
	if _, ok := srw.sz.(*passwdChgSzenario); ok {
		return delay
	}
	srw.retry++
	if srw.retry > 3 {
		srw.retry = 0
		return delay
	}
	cdp.baseHcl.Infof("Szenario %s failed, reschedule faster (retry %d)", srw.sz.Name(), srw.retry)
	delay /= 5
	delay *= time.Duration(srw.retry)
	return delay
}

func (cdp *Engine) reschedule(srw szenarionRunWrapper) {
	if srw.sz.RepeatDelay() < 1 {
		if len(cdp.runChan) < 1 {
			cdp.baseHcl.Infof("No more szenarios, closing the run channel")
			close(cdp.runChan)
		}
		return
	}
	delay := cdp.rescheduleDelay(&srw)
	cdp.baseHcl.Warnf("Rescheduling %s in %v", srw.sz.Name(), delay)
	time.Sleep(delay)
	cdp.runChan <- srw
}

func (cdp *Engine) run() bool {
	ok := true
	if cdp.szenario == nil {
		panic("Szenario must not be nil")
	}
	cdp.hcl.Infof("User %s", cdp.szenario.User().Name())
	cdp.evtMsg = msg.NewSzenarioEvtMsg(cdp.szenario.Name(), cdp.szenario.User().Name(), time.Now())
	engineCancel := cdp.createEngine()
	defer engineCancel()
	timeoutCancel := cdp.TimeOut(cdp.timeout)
	defer timeoutCancel()

	now := time.Now()
	cdp.szenario.User().ResetPasswordIndex()
	defer cdp.reportResults(now) // catches the panic
	if err := cdp.szenario.Execute(cdp); err != nil {
		cdp.hcl.Errorf("Szenario %s returned error: %v", cdp.szenario.Name(), err)
		cdp.ErrorScreenshot(err)
		ok = false
		panic(err)
	}
	return ok
}

func (cdp *Engine) reportResults(start time.Time) {
	d := time.Since(start)

	if cdp.stepDelay > 0 {
		td := d - cdp.stepDelay*time.Duration(len(cdp.stepInfo.stepTimes))
		if td > 0 {
			d = td
		}
	}
	cdp.hcl.Warnf("Running %s took %v", cdp.szenario.Name(), d)
	// cleanup panic, i.e. step failure
	r2 := recover()
	var err error
	if r2 != nil {
		if err2, ok := r2.(error); ok {
			err = err2
		} else {
			err = fmt.Errorf("unknown panic: %v", r2)
		}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		err = fmt.Errorf("timeout reached after %v", d)
	}
	if err != nil {
		cdp.AddErr(err)
	}
	cdp.report(d)
}

// SetStatus sets a status of the event
func (cdp *Engine) SetStatus(key, val string) {
	if cdp.evtMsg == nil {
		cdp.hcl.Warnf("No event not setting status[%s]=%s", key, val)
		return
	}
	cdp.evtMsg.SetStatus(key, val)
}

// AddErr adds a error to the event
func (cdp *Engine) AddErr(err error) {
	if cdp.evtMsg == nil {
		cdp.hcl.Warnf("No event not adding error: %v", err)
		return
	}
	err = fmt.Errorf("%q step %q failed: %w", cdp.szenario.Name(), cdp.stepInfo.name, err)
	cdp.hcl.Errorf("%v", err)
	cdp.evtMsg.AddErr(err)
}

// ErrorScreenshot is used to create a screenshot in case of an error
func (cdp *Engine) ErrorScreenshot(err error) {
	cdp.mu.Lock()
	defer cdp.mu.Unlock()
	name := fmt.Sprintf("error_%s_%s_%s", cdp.szenario.Name(), cdp.stepInfo.name, strcase.ToLowerCamel(err.Error()))
	cdp.hcl.Warnf("Writing failure information: %s", name)
	cdp.dumpHTML(name)
	cdp.screenShot(name)
}

// TimeOut registers a timeout after which the szenario is canceled
func (cdp *Engine) TimeOut(d time.Duration) context.CancelFunc {
	cncl := make(chan any)
	go func() {
		cdp.hcl.Infof("Starting timeout %v", d)
		cdp.mu.Lock()
		cdp.timeoutTicker = time.NewTicker(d)
		cdp.mu.Unlock()
		defer cdp.timeoutTicker.Stop()
		select {
		case <-cdp.timeoutTicker.C:
			cdp.timeoutTicker.Stop()
			cdp.hcl.Warnf("Triggered timeout %v taking screenshot!", d)
			cdp.ErrorScreenshot(fmt.Errorf("timeout %v", d))
			cdp.AddErr(szenario.TimeoutError{Timeout: d})
			time.AfterFunc(500*time.Millisecond, func() {
				if err := chromedp.Cancel(cdp.browser); err != nil {
					cdp.hcl.Warnf("cannot cancel for timeout: %v", err)
				}
			})
		case <-cncl:
			cdp.hcl.Debug("Timeout was canceled")
		case <-cdp.ctx.Done():
		case <-cdp.browser.Done():
		}
	}()
	return func() {
		close(cncl)
	}
}

func (cdp *Engine) addCtxErrs() {
	if cdp.browser.Err() != nil {
		cdp.hcl.Warnf("Browser ctx: %v", cdp.browser.Err())
		cdp.AddErr(fmt.Errorf("Browser ctx: %w", cdp.browser.Err()))
	}
	if cdp.ctx.Err() != nil {
		cdp.hcl.Warnf("Ctx: %v", cdp.ctx.Err())
		cdp.AddErr(fmt.Errorf("Ctx: %w", cdp.ctx.Err()))
	}
}

// ScreenShot saves a screenshot to the outFolder
func (cdp *Engine) ScreenShot(name string) {
	cdp.mu.Lock()
	defer cdp.mu.Unlock()
	cdp.screenShot(name)
}

func (cdp *Engine) screenShot(name string) {
	if cdp.browser.Err() != nil {
		cdp.hcl.Debug("Context is gone not taking screenshot")
		return
	}
	cdp.hcl.Infof("Taking screenshot: %s", name)
	var payload []byte
	if err := chromedp.Run(cdp.browser, chromedp.CaptureScreenshot(&payload)); err != nil {
		cdp.addCtxErrs()
		cdp.hcl.Warnf("cannot get screenshot: %v", err)
		cdp.AddErr(fmt.Errorf("cannot get screenshot: %v", err))
		return
	}
	cdp.evtMsg.AddFile(msg.NewFileMsgItem(
		name,
		mime.Png,
		payload,
	))
}

// DumpHTML saves the HTML to the outFolder
func (cdp *Engine) DumpHTML(name string) {
	cdp.mu.Lock()
	defer cdp.mu.Unlock()
	cdp.dumpHTML(name)
}

func (cdp *Engine) dumpHTML(name string) {
	if cdp.browser.Err() != nil {
		cdp.hcl.Debug("Context is gone not dumping HTML")
		return
	}
	cdp.hcl.Infof("Dumping HTML of %s", name)
	var html string

	err := chromedp.Run(cdp.browser, chromedp.ActionFunc(func(ctx context.Context) error {
		node, err := dom.GetDocument().Do(ctx)
		if err != nil {
			return err
		}
		html, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
		return err
	}))

	if err != nil {
		cdp.addCtxErrs()
		cdp.hcl.Warnf("cannot read dom to get html: %v", err)
		cdp.AddErr(fmt.Errorf("cannot read dom to get html: %v", err))
		return
	}
	cdp.hcl.Debugf("HTML:\n%s", html)
	cdp.evtMsg.AddFile(msg.NewFileMsgItem(
		name,
		mime.HTML,
		[]byte(html),
	))

}
