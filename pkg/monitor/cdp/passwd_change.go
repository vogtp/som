package cdp

import (
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater/user"
)

const (
	PasswordChangeSzenarioName = "Password change"
)

type passwdChgSzenario struct {
	*szenario.Base
	cdp        *Engine
	delay      time.Duration
	szenarios  []szenario.Szenario
	runWrapper szenarionRunWrapper
}

// Execute the szenario
func (s *passwdChgSzenario) Execute(engine szenario.Engine) (err error) {
	hcl := s.cdp.baseHcl
	hcl.Warnf("Running password change")
	for _, sz := range s.szenarios {
		hcl.Warnf("Running password change szenario: %s", sz.Name())
		if err := sz.Execute(engine); err != nil {
			hcl.Errorf("Password change szenarion %q failed: %v", sz.Name(), err)
			engine.AddErr(err)
		}
	}
	hcl.Warnf("Reschedule password chanmge in %v", s.delay)
	time.Sleep(s.delay)
	s.cdp.runChan <- s.runWrapper
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
		Base:  &szenario.Base{},
		cdp:   cdp,
		delay: delay,
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

	pwChgSz.runWrapper = szenarionRunWrapper{sz: pwChgSz, pwChange: true}
	time.Sleep(delay)
	cdp.runChan <- pwChgSz.runWrapper
}
