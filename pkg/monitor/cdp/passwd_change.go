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
	hcl := engine.HCL()
	engine.Step("Inital Check")
	pwCheckInt := 24 * time.Hour
	pwChgCnt := s.User().NumPasswdChg(pwCheckInt)
	hcl.Infof("Number of pw changes: %v in %v", pwChgCnt, pwCheckInt)
	if pwChgCnt >= viper.GetInt(cfg.PasswdChgMax) {
		err := fmt.Errorf("changed %v/%v times in the last %v", pwChgCnt, viper.GetInt(cfg.PasswdChgMax), pwCheckInt)
		hcl.Warnf("Not changing passwords: %v", err)
		return err
	}

	hcl.Warnf("Running password change")
	for _, sz := range s.szenarios {
		hcl.Warnf("Running password change szenario: %s", sz.Name())
		if err := sz.Execute(engine); err != nil {
			hcl.Errorf("Password change szenarion %q failed: %v", sz.Name(), err)
			engine.AddErr(err)
		}
	}
	return nil
}

func (cdp *Engine) passwordChangeLoop(user *user.User) {
	if !viper.GetBool(cfg.PasswdChange) {
		return
	}
	hcl := cdp.baseHcl
	delay := viper.GetDuration(cfg.PasswdChgIntervall)
	hcl.Warnf("Staring password change loop for %s (every %v)", user.Name(), delay)

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
	hcl.Debugf("Found password change szenarios: %v", szNames)
	for _, szName := range szNames {
		sz := szConfig.ByName(szName)
		if sz == nil {
			hcl.Errorf("Password change szenario %s not found", szName)
			continue
		}
		sz.SetUser(user)
		pwChgSz.szenarios = append(pwChgSz.szenarios, sz)
		hcl.Warnf("Added %q as password change szenario.", sz.Name())
	}
	if len(pwChgSz.szenarios) < 1 {
		hcl.Error("No password change szenarios found")
		return
	}

	pwChgSz.runWrapper = szenarionRunWrapper{sz: pwChgSz}
	if viper.GetDuration(cfg.PasswdChangeInitalDelay) > -1 {
		delay = viper.GetDuration(cfg.PasswdChangeInitalDelay)
		hcl.Warnf("Setting initial delay to %v --> ONLY USE THIS IN DEBUGGIN!", delay)
	}
	time.Sleep(delay)
	cdp.runChan <- pwChgSz.runWrapper
}
