package ldapctl

import (
	"fmt"
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
	ldapHost         = "ldap.host"
	ldapBindAccount  = "ldap.bind.account"
	ldapMonitoringOU = "ldap.monitoring.ou"
)

// Command adds all ldap commands
func Command() *cobra.Command {
	flags := ldapCtl.PersistentFlags()
	flags.String(ldapHost, "$host.name$", "FQDN of the ldap host")
	if err := cobra.MarkFlagRequired(flags, ldapHost); err != nil {
		slog.Error("Cannot mark flag as mandatory", "flag", ldapHost, "err", err)
	}
	flags.String(ldapBindAccount, "", "Account to bind to ldap with")
	flags.String(ldapMonitoringOU, "", "OU to use for monitoring stuff")
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

	// name := viper.GetString(ldapBindAccount)
	// if len(name) < 1 {
	// 	return fmt.Errorf("no user given!  Use --%s", ldapBindAccount)
	// }
	// u, err := user.Store.Get(name)
	// if err != nil {
	// 	return fmt.Errorf("cannot get user %s: %v", name, err)
	// }
	result := checks.Result{
		Name:   cmd.Name(),
		Prefix: "ldap",
		Stati:  make(map[string]any),
		CounterFormater: func(name string, value any) string {
			t, ok := value.(time.Duration)
			if !ok {
				return fmt.Sprintf("%v", value)
			}
			return fmt.Sprintf("%vµs", t.Microseconds())
		},
	}
	defer result.PrintExit()
	start := time.Now()
	err := runLdap(&result)
	result.Total = time.Since(start)
	if err != nil {
		result.Err = err
		result.SetCode(icinga.CRITICAL)
	}
	result.Stati[ldapHost] = viper.GetString(ldapHost)
	result.Stati[ldapBindAccount] = viper.GetString(ldapBindAccount)
	result.Stati[ldapMonitoringOU] = viper.GetString(ldapMonitoringOU)
	result.Stati["Version"] = som.Version
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

func runLdap(result *checks.Result) error {
	result.Counter = make(map[string]any)
	start := time.Now()
	defer func() { result.Counter["total"] = time.Since(start) }()
	lc := LDAPClient{
		Host:               viper.GetString(ldapHost),
		Port:               10636,
		UseSSL:             true,
		InsecureSkipVerify: true,
		BindDN:             ldap_user,
		BindPassword:       ladp_password,
	}
	if err := lc.Connect(); err != nil {
		return err
	}
	defer lc.Close()
	result.Counter["connect"] = time.Since(start)
	stepStart := time.Now()
	if err := lc.Bind(); err != nil {
		return err
	}
	result.Counter["bind"] = time.Since(stepStart)
	// slog.Info("Connected", "ldap", lc)

	// ok, _, err := lc.Authenticate(ldap_user, ladp_password)
	// if err != nil {
	// 	t.Fatalf("Cannot auth: %v", err)
	// }

	// if !ok {
	// 	t.Fatal("Not autheticated")
	// }

	testDN := viper.GetString(ldapMonitoringOU)
	if len(testDN) > 0 {
		addTestOU := ldap.NewAddRequest(testDN)
		addTestOU.Attribute("objectClass", []string{"organization"})
		addTestOU.Attribute("objectClass", []string{"top"})
		addTestOU.Attribute("o", []string{"monitoring"})
		addTestOU.Attribute("description", []string{"Container for application monitoring objects"})

		stepStart = time.Now()
		if err := lc.Add(addTestOU); err != nil {
			return err
		}
		result.Counter["add"] = time.Since(stepStart)

		//TODO ldap.NewSearchRequest()

		stepStart = time.Now()
		delTestOU := ldap.NewDelRequest(testDN, nil)
		if err := lc.Del(delTestOU); err != nil {
			return err
		}
		result.Counter["del"] = time.Since(stepStart)
	}
	return nil
}
