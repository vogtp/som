package szenarioctl

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/som/cmd/somctl/term"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/monitor/cdp"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater/user"
)

const (
	// CheckStepBreak adds a breakpont after each step
	CheckStepBreak = "check.step.break"
)

func init() {
	pflag.Bool(CheckStepBreak, false, "adds a breakpont after each step ")
}

var szenarioRun = &cobra.Command{
	Use:     "run",
	Short:   "Run a SOM szenario",
	Long:    `run szenarios, the exit code is equal to the number of failed szenarios`,
	Example: "run all or run owa intranet",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.SetDefault(cfg.BrowserShow, true)
		username := viper.GetString(cfg.CheckUser)
		if len(username) < 1 {
			return errors.New("No user given")
		}
		user, err := user.Store.Get(username)
		if err != nil {
			return fmt.Errorf("could not get user %s is stater connected?: %v", username, err)
		}
		var sz []szenario.Szenario
		szConfig := core.Get().SzenaioConfig()
		ut := user.UserType
		for _, n := range args {
			user.UserType = n
			usz, err := szConfig.ByUser(user)
			if err == nil {
				sz = append(sz, usz...)
				continue
			}
			name := strings.ToLower(n)

			s := szConfig.ByName(name)
			if s != nil {
				sz = append(sz, s)
			} else {
				return fmt.Errorf("no such szenario: %v\nPossible values: %s", n, possibeSzenarioNames())
			}
		}
		user.UserType = ut
		if len(sz) < 1 {
			return fmt.Errorf("no such szenario: %v\nPossible values: %s", strings.Join(args, " "), possibeSzenarioNames())
		}
		if errCnt := runSzenorios(user, sz); errCnt != 0 {
			os.Exit(errCnt)
		}
		return nil
	},
}

func possibeSzenarioNames() string {
	szConfig := core.Get().SzenaioConfig()
	all, _ := szConfig.ByUser(&user.User{UserType: szenario.UserTypeAll})
	names := "\n  User Types: "
	for _, n := range szConfig.GetUserTypes() {
		names = fmt.Sprintf("%s %q", names, n)
	}
	return fmt.Sprintf("%v\n  Szenarios:   %s", names, getNames(all))
}

func getNames(szenarios []szenario.Szenario) string {
	l := ""
	for _, s := range szenarios {
		l = fmt.Sprintf("%s%q ", l, strings.ToLower(s.Name()))
	}
	return l
}

func runSzenorios(user *user.User, szenarios []szenario.Szenario) int {
	slog.Warn("Running szenarios", "szenarios", getNames(szenarios))

	for _, s := range szenarios {
		slog.Info("Starting szenario", log.Szenario, s.Name(), log.User, user.Name())
		s.SetUser(user)
	}
	opts := make([]cdp.Option, 0)
	if viper.GetBool(CheckStepBreak) {
		breackChan := make(chan any)
		defer close(breackChan)
		opts = append(opts, cdp.StepBreakPoint(breackChan))
		go func() {
			for step := range breackChan {
				time.Sleep(150 * time.Millisecond)
				term.Read(fmt.Sprintf("Step %v finished\nAny key to continue...\n", step))
			}
		}()
	}
	cdp, cancel := cdp.New(opts...)
	defer cancel()
	return cdp.Execute(szenarios...)
}
