package cdp

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater/user"
)

const (
	PasswordChangeSzenarioName = "Passwd Change"
)

type passwdChgSzenario struct {
	*szenario.Base
	cdp *Engine
	//	delay      time.Duration
	szenarios  []szenario.Szenario
	runWrapper szenarionRunWrapper
}

// Execute the szenario
func (s *passwdChgSzenario) Execute(engine szenario.Engine) (err error) {
	log := engine.Log()
	engine.Step("Inital Check")
	pwCheckInt := 24 * time.Hour
	pwChgCnt := s.User().NumPasswdChg(pwCheckInt)
	log.Info("Starting password change", "num_pw_chg", pwChgCnt, "interval", pwCheckInt)
	if pwChgCnt >= viper.GetInt(cfg.PasswdChgMax) {
		err := fmt.Errorf("changed %v/%v times in the last %v", pwChgCnt, viper.GetInt(cfg.PasswdChgMax), pwCheckInt)
		log.Warn("Not changing passwords", "reson", err)
		return err
	}

	log.Warn("Running password change")
	for _, sz := range s.szenarios {
		log.Warn("Running password change szenario", "szenario", sz.Name())
		if err := sz.Execute(engine); err != nil {
			log.Error("Password change szenario failed", "szenario", sz.Name(), "error", err)
			engine.AddErr(err)
		}
	}
	return nil
}

func (cdp *Engine) passwordChangeLoop(user *user.User) {
	if !viper.GetBool(cfg.PasswdChange) {
		return
	}
	hcl := cdp.baseLogger
	delay := viper.GetDuration(cfg.PasswdChgIntervall)
	hcl.Warn("Staring password change loop", "szenario", user.Name(), "repeat", delay)

	pwChgSz := &passwdChgSzenario{
		Base: &szenario.Base{
			CheckRepeat: delay,
		},
		cdp: cdp,
	}
	pwChgSz.SetName(PasswordChangeSzenarioName)
	pwChgSz.SetUser(user)

	szConfig := core.Get().SzenaioConfig()
	szNames := viper.GetStringSlice(cfg.PasswdChgSz)
	hcl.Debug("Found password change szenarios", "szenarios", szNames)
	for _, szName := range szNames {
		sz := szConfig.ByName(szName)
		if sz == nil {
			hcl.Error("Password change szenario not found", "szenario", szName)
			continue
		}
		sz.SetUser(user)
		pwChgSz.szenarios = append(pwChgSz.szenarios, sz)
		hcl.Warn("Added password change szenario", "szenario", sz.Name())
	}
	if len(pwChgSz.szenarios) < 1 {
		hcl.Error("No password change szenarios found")
		return
	}

	pwChgSz.runWrapper = szenarionRunWrapper{sz: pwChgSz}
	if viper.GetDuration(cfg.PasswdChangeInitalDelay) > -1 {
		delay = viper.GetDuration(cfg.PasswdChangeInitalDelay)
		hcl.Warn("Setting initial delay --> ONLY USE THIS IN DEBUGGIN!", "delay", delay)
	}
	time.Sleep(delay)
	cdp.runChan <- pwChgSz.runWrapper
}
