package userctl

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/cmd/somctl/term"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater/user"
)

// Command adds all user commands
func Command() *cobra.Command {
	userCtl.AddCommand(userList)
	userCtl.AddCommand(userShow)
	userCtl.AddCommand(userAdd)
	if hcl.IsGoRun() {
		// show password is only supported in debugging
		userShow.AddCommand(userShowPw)
	}
	return userCtl
}

var userCtl = &cobra.Command{
	Use:   "user",
	Short: "Manage SOM users",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var userShow = &cobra.Command{
	Use:   "show USERNAME",
	Short: "Display a SOM user",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		u, err := user.Store.Get(name)
		if err != nil {
			return fmt.Errorf("cannot get user %s: %v", name, err)
		}
		fmt.Printf("User %s:\n%v\n", name, u.String())
		return nil
	},
}

var userAdd = &cobra.Command{
	Use:     "add [USERNAME EMAIL TYPE PASSWORD]",
	Short:   "Add or modify a SOM user",
	Long:    ``,
	Aliases: []string{"modify", "mod", "set"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("\nAdd new user:")
		defaultValue := ""
		name := term.ReadOrArgs("username", args, 0, defaultValue)
		existUser, err := user.Store.Get(name)
		if err != nil {
			core.Get().HCL().Debugf("User %s not found in backend: %v\n", name, err)
		}
		if existUser != nil {
			fmt.Printf("Found existing user %s empty input will be filled with values from this user...\n", existUser.Name())
			defaultValue = existUser.Mail
		}
		email := term.ReadOrArgs("email", args, 1, defaultValue)
		if existUser != nil {
			defaultValue = existUser.Type()
		}
		ty := ""
		szConfig := core.Get().SzenaioConfig()
		for len(ty) < 1 {
			ty = term.ReadOrArgs("type", args, 2, defaultValue)
			ty = strings.TrimSpace(ty)
			if szConfig != nil && len(szConfig.GetUserTypes()) > 1 {
				ut := szConfig.GetUserType(ty)
				if ut == nil {
					fmt.Printf("%q is not a valid user type\nPossible types are:\n", ty)
					for _, k := range szConfig.GetUserTypes() {
						if k == szenario.UserTypeAll {
							continue
						}
						fmt.Printf("  %s: %s\n", k, szConfig.GetUserType(k).Desc)
					}
					ty = ""
				}
			}
		}
		pw := term.Password("password", args, 3)
		u := &user.User{
			Username: name,
			Mail:     email,
			UserType: ty,
		}
		if existUser != nil {
			u.History = existUser.History
		}
		u.SetPassword(pw)
		if err := u.IsValid(); err != nil {
			return fmt.Errorf("user is not valid: %w", err)
		}
		err = user.Store.Save(u)
		if err != nil {
			return fmt.Errorf("cannot add user %s: %v", name, err)
		}
		fmt.Printf("Added %s:\n%v\n", name, u.String())
		return nil
	},
}

var userShowPw = &cobra.Command{
	Use:   "pass USERNAME",
	Short: "Display a SOM users password",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		u, err := user.Store.Get(name)
		if err != nil {
			return fmt.Errorf("cannot set password of user %s: %v", name, err)
		}
		fmt.Printf("\n%s: %s\n", name, u.String())
		fmt.Printf("  %-20s (Created: %s LastUse: %s)\n", u.Password(), u.PasswordCreated().Format(cfg.TimeFormatString), u.PasswordLastUse().Format(cfg.TimeFormatString))
		fmt.Println("History:")
		for pw := u.NextPassword(); pw != ""; pw = u.NextPassword() {
			fmt.Printf("  %-20s (Created: %s LastUse: %s)\n", pw, u.PasswordCreated().Format(cfg.TimeFormatString), u.PasswordLastUse().Format(cfg.TimeFormatString))
		}
		return nil
	},
}

var userList = &cobra.Command{
	Use:   "list",
	Short: "List SOM users",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("\nUsers:")
		users, err := user.Store.List()
		if err != nil {
			return fmt.Errorf("cannot get userlist: %v", err)
		}
		sort.Slice(users, func(i, j int) bool {
			return users[i].Name() < users[j].Name()
		})
		for i, u := range users {
			fmt.Printf("%v: %s\n", i, u.String())
		}
		return nil
	},
}
