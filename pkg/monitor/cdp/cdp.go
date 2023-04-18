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

	// runChan contains the next szenarios to be run
	runChan chan szenarionRunWrapper
	// triggerSzenarioChan triggers a szenario run (send the szenario name to the channel or all for all szenarios)
	triggerSzenarioChan chan string
	// stepBreakPoint is used to pause between steps, read from chan to proceed to next step
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

	timeout time.Duration

	// flags
	headless      bool
	noClose       bool
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
		ctx:                 ctx,
		baseHcl:             hcl,
		hcl:                 hcl,
		bus:                 core.Bus(),
		runChan:             make(chan szenarionRunWrapper, 100),
		triggerSzenarioChan: make(chan string, 100),
		sendReport:          true,
		consMsg:             make(map[string]int),
		headless:            !viper.GetBool(cfg.BrowserShow),
		noClose:             viper.GetBool(cfg.BrowserNoClose),
		stepDelay:           viper.GetDuration(cfg.CheckStepDelay),
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
		cdp.hcl.Info("Updating timeout", "timeout", d)
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
	cdp.hcl.Info("Scheduling szenarios", "count", len(szenarios))
	for _, sz := range szenarios {
		cdp.hcl.Info("Scheduling szenario", "szenario", sz.Name())
		cdp.runChan <- szenarionRunWrapper{sz: sz}
	}
}

// loop runs the szenarios and never returns
func (cdp *Engine) loop() {
	for srw := range cdp.runChan {
		cdp.hcl = cdp.baseHcl.With("szenario", srw.sz.Name())
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
	delay /= 5
	delay *= time.Duration(srw.retry)
	cdp.baseHcl.Info("Szenario failed, reschedule faster", "szenario", srw.sz.Name(), "retry", srw.retry, "delay", delay)
	return delay
}

func (cdp *Engine) reschedule(srw szenarionRunWrapper) {
	if srw.sz.RepeatDelay() < 1 {
		if len(cdp.runChan) < 1 {
			cdp.baseHcl.Info("No more szenarios, closing the run channel")
			close(cdp.runChan)
		}
		return
	}
	delay := cdp.rescheduleDelay(&srw)
	cdp.baseHcl.Warn("Rescheduling", "szenario", srw.sz.Name(), "delay", delay)

	ticker := time.NewTicker(delay)
	for wait := true; wait; {
		select {
		case <-ticker.C:
			wait = false
		case sz := <-cdp.triggerSzenarioChan:
			switch sz {
			case srw.sz.Name():
				wait = false
			case "all":
				wait = false
			}
		}
	}
	cdp.runChan <- srw
}

func (cdp *Engine) run() bool {
	ok := true
	if cdp.szenario == nil {
		panic("Szenario must not be nil")
	}
	cdp.hcl = cdp.hcl.With("user", cdp.szenario.User().Name())
	cdp.evtMsg = msg.NewSzenarioEvtMsg(cdp.szenario.Name(), cdp.szenario.User().Name(), time.Now())
	engineCancel := cdp.createEngine()
	defer engineCancel()
	timeoutCancel := cdp.TimeOut(cdp.szenario.Timeout())
	defer timeoutCancel()

	now := time.Now()
	cdp.szenario.User().ResetPasswordIndex()
	defer func() {
		if err := cdp.szenario.User().Save(); err != nil {
			cdp.hcl.Error("Cannot save user after szenario run", "error", err)
		}
	}()
	defer cdp.reportResults(now) // catches the panic
	if err := cdp.szenario.Execute(cdp); err != nil {
		cdp.hcl.Error("Szenario returned error", "error", err)
		cdp.ErrorScreenshot(err)
		cdp.evtMsg.SetStatus("URL", cdp.GetURL())
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
	cdp.hcl.Warn("Szenario finished", "duration", d)
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
		cdp.hcl.Warn("No event not setting status", "key", key, "value", val)
		return
	}
	cdp.evtMsg.SetStatus(key, val)
}

// AddErr adds a error to the event
func (cdp *Engine) AddErr(err error) {
	if cdp.evtMsg == nil {
		cdp.hcl.Warn("No event not adding error", "error", err)
		return
	}
	err = fmt.Errorf("%q step %q failed: %w", cdp.szenario.Name(), cdp.stepInfo.name, err)
	cdp.hcl.Warn("Adding error to event", "error", err)
	cdp.evtMsg.AddErr(err)
}

// ErrorScreenshot is used to create a screenshot in case of an error
func (cdp *Engine) ErrorScreenshot(err error) {
	cdp.mu.Lock()
	defer cdp.mu.Unlock()
	name := fmt.Sprintf("error_%s_%s_%s", cdp.szenario.Name(), cdp.stepInfo.name, strcase.ToLowerCamel(err.Error()))
	cdp.hcl.Warn("Writing failure information", "screenshot", name)
	cdp.dumpHTML(name)
	cdp.screenShot(name)
}

// TimeOut registers a timeout after which the szenario is canceled
func (cdp *Engine) TimeOut(d time.Duration) context.CancelFunc {
	cncl := make(chan any)
	go func() {
		cdp.hcl.Info("Starting timeout", "timeout", d)
		cdp.mu.Lock()
		cdp.timeoutTicker = time.NewTicker(d)
		cdp.mu.Unlock()
		defer cdp.timeoutTicker.Stop()
		select {
		case <-cdp.timeoutTicker.C:
			cdp.timeoutTicker.Stop()
			cdp.hcl.Warn("Triggered timeout taking screenshot!", "timeout", d)
			cdp.ErrorScreenshot(fmt.Errorf("timeout %v", d))
			cdp.AddErr(szenario.TimeoutError{Timeout: d})
			cdp.evtMsg.SetStatus("URL", cdp.GetURL())
			time.AfterFunc(500*time.Millisecond, func() {
				if err := chromedp.Cancel(cdp.browser); err != nil {
					cdp.hcl.Warn("cannot cancel for timeout", "error", err)
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
		cdp.hcl.Warn("Browser ctx error", "error", cdp.browser.Err())
		cdp.AddErr(fmt.Errorf("Browser ctx: %w", cdp.browser.Err()))
	}
	if cdp.ctx.Err() != nil {
		cdp.hcl.Warn("Ctx error", "error", cdp.ctx.Err())
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
	cdp.hcl.Info("Taking screenshot", "screenshot", name)
	var payload []byte
	if err := chromedp.Run(cdp.browser, chromedp.CaptureScreenshot(&payload)); err != nil {
		cdp.addCtxErrs()
		cdp.hcl.Warn("cannot get screenshot", "error", err)
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
	cdp.hcl.Info("Dumping HTML", "szenario", name)
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
		cdp.hcl.Warn("cannot read dom to get html", "error", err)
		cdp.AddErr(fmt.Errorf("cannot read dom to get html: %v", err))
		return
	}
	cdp.hcl.Debug("HTML:", "html", html)
	cdp.evtMsg.AddFile(msg.NewFileMsgItem(
		name,
		mime.HTML,
		[]byte(html),
	))

}
