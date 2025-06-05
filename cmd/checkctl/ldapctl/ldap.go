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
	result := checks.NewCheckResult(cmd.Name(), checks.CheckPrefix("ldap"), checks.CounterFormater(timeFormater))
	defer result.PrintExit()
	start := time.Now()
	err := runLdap(result)
	result.Total = time.Since(start)
	if err != nil {
		result.SetError(err)
	}
	result.SetStatus(ldapHost, viper.GetString(ldapHost))
	result.SetStatus(ldapBindAccount, viper.GetString(ldapBindAccount))
	result.SetStatus(ldapMonitoringOU, viper.GetString(ldapMonitoringOU))
	result.SetStatus("Version", som.Version)
	return err
}

func timeFormater(name string, value any) string {
	t, ok := value.(time.Duration)
	if !ok {
		return fmt.Sprintf("%v", value)
	}
	return fmt.Sprintf("%vus", t.Microseconds())
}

func runLdap(result *checks.Result) error {
	start := time.Now()
	defer func() { result.SetCounter("total", time.Since(start)) }()
	lc := &LDAPClient{
		Host:               viper.GetString(ldapHost),
		Port:               10636,
		UseSSL:             true,
		InsecureSkipVerify: true,
		BindDN:             ldap_user,
		BindPassword:       ladp_password,
	}
	if err := lc.Connect(); err != nil {
		result.SetCode(icinga.CRITICAL)
		return err
	}
	defer lc.Close()
	result.SetCounter("connect", time.Since(start))
	stepStart := time.Now()
	if err := lc.Bind(); err != nil {
		result.SetCode(icinga.CRITICAL)
		return err
	}
	result.SetCounter("bind", time.Since(stepStart))

	//TODO ldap.NewSearchRequest()

	if err := createAndDeleteOU(lc, result); err != nil {
		result.SetError(err)
	}

	return nil
}

func createAndDeleteOU(lc *LDAPClient, result *checks.Result) error {
	testDN := viper.GetString(ldapMonitoringOU)
	if len(testDN) < -1 {
		return nil
	}
	addTestOU := ldap.NewAddRequest(testDN)
	addTestOU.Attribute("objectClass", []string{"organization"})
	addTestOU.Attribute("objectClass", []string{"top"})
	addTestOU.Attribute("o", []string{"monitoring"})
	addTestOU.Attribute("description", []string{"Container for application monitoring objects"})

	stepStart := time.Now()
	if err := lc.Add(addTestOU); err != nil {
		return err
	}
	result.SetCounter("add", time.Since(stepStart))

	stepStart = time.Now()
	delTestOU := ldap.NewDelRequest(testDN, nil)
	if err := lc.Del(delTestOU); err != nil {
		return err
	}
	result.SetCounter("del", time.Since(stepStart))
	return nil
}
