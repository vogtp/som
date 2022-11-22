package szenario

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/cfg"
)

type setNamer interface {
	SetName(string)
}

// Add a szenario to the engine
func (c Config) Add(name string, s Szenario, ut []*UserType) Szenario {
	if sn, ok := s.(setNamer); ok {
		sn.SetName(name)
	}
	if len(s.Name()) < 1 {
		panic(fmt.Errorf("Szenaio (%T) must have a name", s))
	}
	hcl.Tracef("Initialising szenario %s", s.Name())
	for _, sz := range c.userTypes[c.allSz.Name].Szenarios {
		if sz.Name() == s.Name() {
			hcl.Errorf("Szenario %s already exists", s.Name())
			//return sz
		}
	}
	if err := c.addUserType(c.allSz, s); err != nil {
		hcl.Errorf("Cannot add szenario %s to usertype all: %v", s.Name(), err)
	}
	for _, t := range ut {
		if err := c.addUserType(t, s); err != nil {
			hcl.Errorf("Adding usertype %s to szenario %s: %v", t, s.Name(), err)
		}
	}
	return s
}

// Base is the base type of all szenarios
type Base struct {
	name         string
	user         User
	CheckRepeat  time.Duration
	CheckTimeout time.Duration
	LoginRetry   int
}

// SetName do not call!
// rename panics
func (s *Base) SetName(name string) {
	if len(s.name) > 0 {
		panic("Szenario renaming is not supported")
	}
	s.name = name
}

// Name returns the name of the szenario
func (s Base) Name() string {
	return s.name
}

// User returns the user the szenario runs with
func (s *Base) User() User {
	return s.user
}

// SetUser set the user the szenario runs with
func (s *Base) SetUser(u User) {
	s.user = u
}

// MaxLoginTry returns how many times a login with a new password should be attemped
func (s Base) MaxLoginTry() int {
	if s.LoginRetry > 0 {
		return s.LoginRetry
	}
	return 4
}

// RepeatDelay between executions
func (s Base) RepeatDelay() time.Duration {
	if s.CheckRepeat > 0 {
		return s.CheckRepeat
	}
	return viper.GetDuration(cfg.CheckRepeat)
}

// Timeout for an execution
func (s Base) Timeout() time.Duration {
	if s.CheckTimeout > 0 {
		return s.CheckTimeout
	}
	return viper.GetDuration(cfg.CheckTimeout)
}
