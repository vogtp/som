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
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/mime"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater/user"
	"golang.org/x/exp/slog"
)

// Option configures the Engine
type Option func(*Engine)

// Engine contains all information to run a chrome szenario
type Engine struct {
	ctx        context.Context
	baseLogger *slog.Logger
	log        *slog.Logger
	bus        *core.Bus

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
	slog := core.Log().With(log.Component, "cdp")
	ctx, cancel := context.WithCancel(context.Background())
	cdp := &Engine{
		ctx:                 ctx,
		baseLogger:          slog,
		log:                 slog,
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
		cdp.log.Warn("Canceled by OS signal")
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
		cdp.log.Info("Updating timeout", "timeout", d)
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
	cdp.log.Info("Scheduling szenarios", "count", len(szenarios))
	for _, sz := range szenarios {
		cdp.log.Info("Scheduling szenario", log.Szenario, sz.Name())
		cdp.runChan <- szenarionRunWrapper{sz: sz}
	}
}

// loop runs the szenarios and never returns
func (cdp *Engine) loop() {
	for srw := range cdp.runChan {
		cdp.log = cdp.baseLogger.With(log.Szenario, srw.sz.Name())
		cdp.szenario = srw.sz
		cdp.stepInfo = newStepInfo(cdp.log)
		srw.lastRunOk = cdp.run()
		cdp.log = cdp.baseLogger
		go cdp.reschedule(srw)
	}
	cdp.log.Info("Finished main szenario loop")
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
	cdp.baseLogger.Info("Szenario failed, reschedule faster", log.Szenario, srw.sz.Name(), "retry", srw.retry, "delay", delay)
	return delay
}

func (cdp *Engine) reschedule(srw szenarionRunWrapper) {
	if srw.sz.RepeatDelay() < 1 {
		if len(cdp.runChan) < 1 {
			cdp.baseLogger.Info("No more szenarios, closing the run channel")
			close(cdp.runChan)
		}
		return
	}
	delay := cdp.rescheduleDelay(&srw)
	cdp.baseLogger.Warn("Rescheduling", log.Szenario, srw.sz.Name(), "delay", delay)

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
	cdp.log = cdp.log.With(log.User, cdp.szenario.User().Name())
	cdp.evtMsg = msg.NewSzenarioEvtMsg(cdp.szenario.Name(), cdp.szenario.User().Name(), time.Now())
	engineCancel := cdp.createEngine()
	defer engineCancel()
	timeoutCancel := cdp.TimeOut(cdp.szenario.Timeout())
	defer timeoutCancel()

	now := time.Now()
	cdp.szenario.User().ResetPasswordIndex()
	defer func() {
		if err := cdp.szenario.User().Save(); err != nil {
			cdp.log.Error("Cannot save user after szenario run", log.Error, err)
		}
	}()
	defer cdp.reportResults(now) // catches the panic
	if err := cdp.szenario.Execute(cdp); err != nil {
		cdp.log.Error("Szenario returned error", log.Error, err)
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
	cdp.log.Warn("Szenario finished", "duration", d)
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
		cdp.log.Warn("No event not setting status", "key", key, "value", val)
		return
	}
	cdp.evtMsg.SetStatus(key, val)
}

// AddErr adds a error to the event
func (cdp *Engine) AddErr(err error) {
	if cdp.evtMsg == nil {
		cdp.log.Warn("No event not adding error", log.Error, err)
		return
	}
	err = fmt.Errorf("%q step %q failed: %w", cdp.szenario.Name(), cdp.stepInfo.name, err)
	cdp.log.Warn("Adding error to event", log.Error, err)
	cdp.evtMsg.AddErr(err)
}

// ErrorScreenshot is used to create a screenshot in case of an error
func (cdp *Engine) ErrorScreenshot(err error) {
	cdp.mu.Lock()
	defer cdp.mu.Unlock()
	name := fmt.Sprintf("error_%s_%s_%s", cdp.szenario.Name(), cdp.stepInfo.name, strcase.ToLowerCamel(err.Error()))
	cdp.log.Warn("Writing failure information", "screenshot", name)
	cdp.dumpHTML(name)
	cdp.screenShot(name)
}

// TimeOut registers a timeout after which the szenario is canceled
func (cdp *Engine) TimeOut(d time.Duration) context.CancelFunc {
	cncl := make(chan any)
	go func() {
		cdp.log.Info("Starting timeout", "timeout", d)
		cdp.mu.Lock()
		cdp.timeoutTicker = time.NewTicker(d)
		cdp.mu.Unlock()
		defer cdp.timeoutTicker.Stop()
		select {
		case <-cdp.timeoutTicker.C:
			cdp.timeoutTicker.Stop()
			cdp.log.Warn("Triggered timeout taking screenshot!", "timeout", d)
			cdp.ErrorScreenshot(fmt.Errorf("timeout %v", d))
			cdp.AddErr(szenario.TimeoutError{Timeout: d})
			cdp.evtMsg.SetStatus("URL", cdp.GetURL())
			time.AfterFunc(500*time.Millisecond, func() {
				if err := chromedp.Cancel(cdp.browser); err != nil {
					cdp.log.Warn("cannot cancel for timeout", log.Error, err)
				}
			})
		case <-cncl:
			cdp.log.Debug("Timeout was canceled")
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
		cdp.log.Warn("Browser ctx error", log.Error, cdp.browser.Err())
		cdp.AddErr(fmt.Errorf("Browser ctx: %w", cdp.browser.Err()))
	}
	if cdp.ctx.Err() != nil {
		cdp.log.Warn("Ctx error", log.Error, cdp.ctx.Err())
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
		cdp.log.Debug("Context is gone not taking screenshot")
		return
	}
	cdp.log.Info("Taking screenshot", "screenshot", name)
	var payload []byte
	if err := chromedp.Run(cdp.browser, chromedp.CaptureScreenshot(&payload)); err != nil {
		cdp.addCtxErrs()
		cdp.log.Warn("cannot get screenshot", log.Error, err)
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
		cdp.log.Debug("Context is gone not dumping HTML")
		return
	}
	cdp.log.Info("Dumping HTML", log.Szenario, name)
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
		cdp.log.Warn("cannot read dom to get html", log.Error, err)
		cdp.AddErr(fmt.Errorf("cannot read dom to get html: %v", err))
		return
	}
	cdp.log.Debug("HTML:", "html", html)
	cdp.evtMsg.AddFile(msg.NewFileMsgItem(
		name,
		mime.HTML,
		[]byte(html),
	))

}
