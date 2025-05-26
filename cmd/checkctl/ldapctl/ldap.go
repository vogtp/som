package ldapctl

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/ldap.v2"
)

// Command adds all ldap commands
func Command() *cobra.Command {
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
	timing, err := runLdap()
	total := timing["total"].Milliseconds()
	ret := fmt.Sprintf("%s OK - duration %vms", strings.ToUpper(cmd.Name()), total)
	if err != nil {
		ret = fmt.Sprintf("%s CRITICAL - %v", strings.ToUpper(cmd.Name()), err)
	}
	pref := ""
	disp := ""
	for n, t := range timing {
		pref = fmt.Sprintf("%s%s_ms=%v ", pref, n, t.Milliseconds())
		pref = fmt.Sprintf("%s%s=%vms ", pref, n, t.Milliseconds())
		disp = fmt.Sprintf("%s%s\t%vms\n", disp, n, t.Milliseconds())
	}
	// disp = fmt.Sprintf("%s%s\t%vms", disp, "total", total)
	fmt.Printf("%s\n\n%s | %s", ret, disp, pref)
	if err != nil {
		os.Exit(2)
	}
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

	delTestOU := ldap.NewDelRequest(testDN, nil)
	stepStart = time.Now()
	if err = lc.Del(delTestOU); err != nil {
		return
	}
	timing["del"] = time.Since(stepStart)
	return
}
