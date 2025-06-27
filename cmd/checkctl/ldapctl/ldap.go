package ldapctl

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/check"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/som"
	"gopkg.in/ldap.v2"
)

const (
	ldapHost         = "ldap.host"
	ldapBindDN       = "ldap.bind.DN"
	ldapSearchFilter = "ldap.search.filter"
	ldapMonitoringOU = "ldap.monitoring.ou"
)

// Command adds all ldap commands
func Command() *cobra.Command {
	flags := ldapCtl.PersistentFlags()
	flags.String(ldapHost, "$host.name$", "FQDN of the ldap host")
	flags.String(ldapBindDN, "uid=admin", "Account DN to bind to ldap with")
	flags.String(ldapSearchFilter, "(uid=vogtp)", "Searchfilter for the LDAP search")
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		check.SetWarningThresholdDefault("connect:20ms total:100ms")
		check.SetCriticalThresholdDefault("connect:60ms total:300ms")
		return nil
	},
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
	result := check.NewResult(cmd.Name(), check.CounterFormater(timeFormater))
	defer result.PrintExit()
	start := time.Now()
	err := runLdap(result)
	result.SetHeader("%s", fmt.Sprintf("Duration %s", timeFormater("total", check.Data{Value: time.Since(start)})))
	if err != nil {
		result.SetError(err)
	}
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if !strings.HasPrefix(f.Name, "ldap") {
			return
		}
		result.SetStatus(f.Name, f.Value)
	})
	result.SetStatus("Version", som.Version)
	return err
}

func timeFormater(name string, value check.Data) string {
	t, ok := value.Value.(time.Duration)
	if !ok {
		return fmt.Sprintf("%v", value)
	}
	return fmt.Sprintf("%vus", t.Microseconds())
}

func runLdap(result *check.Result) error {
	start := time.Now()
	defer func() { result.SetCounter("total", time.Since(start)) }()
	lc := &LDAPClient{
		Host:               viper.GetString(ldapHost),
		Port:               10636,
		UseSSL:             true,
		InsecureSkipVerify: true,
		BindDN:             viper.GetString(ldapBindDN),
		BindPassword:       ladp_password,
	}
	slog := slog.With(ldapHost, lc.Host, ldapBindDN, lc.BindDN)
	slog.Debug("Connecting to LDAP host")
	if err := lc.Connect(); err != nil {
		slog.Warn("Cannot connect to LDAP host", "host", lc.Host, "bindDB", lc.BindDN)
		result.SetCode(icinga.CRITICAL)
		return err
	}
	defer lc.Close()
	result.SetCounter("connect", time.Since(start))
	// if not user is give just check the connect
	if len(lc.BindDN) < 1 {
		return nil
	}
	stepStart := time.Now()
	if err := lc.Bind(); err != nil {
		result.SetCode(icinga.CRITICAL)
		return err
	}
	result.SetCounter("bind", time.Since(stepStart))

	if err := search(lc, result); err != nil {
		result.SetError(fmt.Errorf("LDAP Search: %w", err))
	}

	if err := createAndDeleteOU(lc, result); err != nil {
		result.SetError(fmt.Errorf("Create/Delete OU: %w", err))
	}

	return nil
}

func search(l *LDAPClient, result *check.Result) error {
	filter := viper.GetString(ldapSearchFilter)
	if len(filter) < 1 {
		return nil
	}
	searchRequest := ldap.NewSearchRequest(
		"",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		// fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"dn"},
		nil,
	)

	stepStart := time.Now()
	sr, err := l.Search(searchRequest)
	if err != nil {
		return fmt.Errorf("ldap search: %w", err)
	}
	result.SetCounter("search", time.Since(stepStart))

	if len(sr.Entries) < 1 {
		return fmt.Errorf("User does not exist")
	}
	slog.Info("LDAP search result", "filter", filter, "result", sr, "fist DN", sr.Entries[0].DN)

	return nil
}

func createAndDeleteOU(lc *LDAPClient, result *check.Result) error {
	testDN := viper.GetString(ldapMonitoringOU)
	if len(testDN) < 1 {
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
