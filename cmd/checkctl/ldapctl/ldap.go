package ldapctl

import (
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/checks"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/som"
	"gopkg.in/ldap.v2"
)

const (
	ldapHost        = "ldap.host"
	ldapBindAccount = "ldap.bind.account"
)

// Command adds all ldap commands
func Command() *cobra.Command {
	flags := ldapCtl.PersistentFlags()
	flags.String(ldapHost, "$host.name$", "FQDN of the ldap host")
	if err := cobra.MarkFlagRequired(flags, ldapHost); err != nil {
		slog.Error("Cannot mark flag as mandatory", "flag", ldapHost, "err", err)
	}
	flags.String(ldapBindAccount, "", "Account to bind to ldap with")
	flags.VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})
	return ldapCtl
}

var ldapCtl = &cobra.Command{
	Use:   "ldap",
	Short: "Check a ldap",
	Long:  ``,

	RunE: func(cmd *cobra.Command, args []string) error {
		return LdapCheckCmd(cmd, args)
	},
}

func LdapCheckCmd(cmd *cobra.Command, args []string) error {

	// name := viper.GetString(cfg.CheckUser)
	// if len(name) < 1 {
	// 	return fmt.Errorf("no user given!  Use --%s", cfg.CheckUser)
	// }
	// u, err := user.Store.Get(name)
	// if err != nil {
	// 	return fmt.Errorf("cannot get user %s: %v", name, err)
	// }
	result := checks.Result{
		Name:   cmd.Name(),
		Prefix: "ldap",
		Result: icinga.OK,
		Stati:  make(map[string]any),
	}
	timing, err := runLdap()
	result.Total = timing["total"]
	result.Timing = timing
	if err != nil {
		result.Err = err
		result.Result = icinga.CRITICAL
	}
	result.Stati[ldapHost] = viper.GetString(ldapHost)
	result.Stati[ldapBindAccount] = viper.GetString(ldapBindAccount)
	result.Stati["Version"] = som.Version
	result.PrintExit()
	return err
}

/*
https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/pluginapi.html
Plugin Return Code	Service State	Host State
0	OK	UP
1	WARNING	UP or DOWN/UNREACHABLE*
2	CRITICAL	DOWN/UNREACHABLE
3	UNKNOWN	DOWN/UNREACHABLE
*/

func runLdap() (timing map[string]time.Duration, err error) {
	timing = make(map[string]time.Duration)
	start := time.Now()
	defer func() { timing["total"] = time.Since(start) }()
	lc := LDAPClient{
		//BindDN: "uid=admin",
		Host:               "its-ds-ngi-dev-1.its.unibas.ch",
		Port:               10636,
		UseSSL:             true,
		InsecureSkipVerify: true,
		BindDN:             ldap_user,
		BindPassword:       ladp_password,
	}
	if err = lc.Connect(); err != nil {
		return
	}
	defer lc.Close()
	timing["connect"] = time.Since(start)
	stepStart := time.Now()
	if err = lc.Bind(); err != nil {
		return
	}
	timing["bind"] = time.Since(stepStart)
	// slog.Info("Connected", "ldap", lc)

	// ok, _, err := lc.Authenticate(ldap_user, ladp_password)
	// if err != nil {
	// 	t.Fatalf("Cannot auth: %v", err)
	// }

	// if !ok {
	// 	t.Fatal("Not autheticated")
	// }

	testDN := "o=monitoring,dc=unibas,dc=ch"
	addTestOU := ldap.NewAddRequest(testDN)
	addTestOU.Attribute("objectClass", []string{"organization"})
	addTestOU.Attribute("objectClass", []string{"top"})
	addTestOU.Attribute("o", []string{"monitoring"})
	addTestOU.Attribute("description", []string{"Container for application monitoring objects"})

	stepStart = time.Now()
	if err = lc.Add(addTestOU); err != nil {
		return
	}
	timing["add"] = time.Since(stepStart)

	//TODO ldap.NewSearchRequest()

	stepStart = time.Now()
	delTestOU := ldap.NewDelRequest(testDN, nil)
	if err = lc.Del(delTestOU); err != nil {
		return
	}
	timing["del"] = time.Since(stepStart)
	return
}
