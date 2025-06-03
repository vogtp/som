// Package ldapctl provides a simple ldap client to authenticate,
// retrieve basic information and groups for a user.
package ldapctl

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gopkg.in/ldap.v2"
)

type LDAPClient struct {
	*ldap.Conn

	Attributes         []string
	Base               string
	BindDN             string
	BindPassword       string
	GroupFilter        string // e.g. "(memberUid=%s)"
	Host               string
	ServerName         string
	UserFilter         string // e.g. "(uid=%s)"
	Port               int
	InsecureSkipVerify bool
	UseSSL             bool
	SkipTLS            bool
	ClientCertificates []tls.Certificate // Adding client certificates
}

// Connect connects to the ldap backend.
func (lc *LDAPClient) Connect() error {
	ldap.DefaultTimeout = time.Second
	if lc.Conn == nil {
		var l *ldap.Conn
		var err error
		address := fmt.Sprintf("%s:%d", lc.Host, lc.Port)
		if !lc.UseSSL {
			l, err = ldap.Dial("tcp", address)
			if err != nil {
				return err
			}

			// Reconnect with TLS
			if !lc.SkipTLS {
				err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
				if err != nil {
					return err
				}
			}
		} else {
			config := &tls.Config{
				InsecureSkipVerify: lc.InsecureSkipVerify,
				ServerName:         lc.ServerName,
			}
			if len(lc.ClientCertificates) > 0 {
				config.Certificates = lc.ClientCertificates
			}
			l, err = ldap.DialTLS("tcp", address, config)
			if err != nil {
				return err
			}
		}

		lc.Conn = l
	}
	return nil
}

var _ slog.LogValuer = (*LDAPClient)(nil)

func (lc LDAPClient) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("host", lc.Host),
		slog.Int("port", lc.Port),
		slog.String("bind_dn", lc.BindDN),
	)
}

// Close closes the ldap backend connection.
func (lc *LDAPClient) Close() {
	if lc.Conn != nil {
		lc.Conn.Close()
		lc.Conn = nil
	}
}

func (lc *LDAPClient) Bind() error {
	err := lc.Connect()
	if err != nil {
		return err
	}

	if len(lc.BindDN) < 1 {
		return fmt.Errorf("no BindDN given")
	}
	if len(lc.BindPassword) < 1 {
		return fmt.Errorf("no BindPassword given")
	}
	return lc.Conn.Bind(lc.BindDN, lc.BindPassword)
}

// Authenticate authenticates the user against the ldap backend.
func (lc *LDAPClient) Authenticate(username, password string) (bool, map[string]string, error) {
	if err := lc.Bind(); err != nil {
		return false, nil, err
	}
	if len(lc.UserFilter) < 1 {
		lc.UserFilter = "(uid=%s)"
	}
	if len(lc.GroupFilter) < 1 {
		lc.GroupFilter = "(i9itl9dCETofzYXQsTwA=%s)"
	}
	attributes := append(lc.Attributes, "dn")
	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.UserFilter, username),
		attributes,
		nil,
	)

	sr, err := lc.Conn.Search(searchRequest)
	if err != nil {
		return false, nil, err
	}

	if len(sr.Entries) < 1 {
		return false, nil, errors.New("User does not exist")
	}

	if len(sr.Entries) > 1 {
		return false, nil, errors.New("Too many entries returned")
	}

	for i, e := range sr.Entries {
		fmt.Printf("%v: %v\n", i, e)
	}

	userDN := sr.Entries[0].DN
	user := map[string]string{}
	for _, attr := range lc.Attributes {
		user[attr] = sr.Entries[0].GetAttributeValue(attr)
	}

	// Bind as the user to verify their password
	err = lc.Conn.Bind(userDN, password)
	if err != nil {
		return false, user, err
	}

	// Rebind as the read only user for any further queries
	if lc.BindDN != "" && lc.BindPassword != "" {
		err = lc.Conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return true, user, err
		}
	}

	return true, user, nil
}

// GetGroupsOfUser returns the group for a user.
func (lc *LDAPClient) GetGroupsOfUser(username string) ([]string, error) {
	err := lc.Connect()
	if err != nil {
		return nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.GroupFilter, username),
		[]string{"cn"}, // can it be something else than "cn"?
		nil,
	)
	sr, err := lc.Conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	groups := []string{}
	for _, entry := range sr.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}
	return groups, nil
}
