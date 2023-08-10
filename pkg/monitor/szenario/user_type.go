package szenario

import (
	"fmt"
	"strings"

	"log/slog"
)

// MustUserType panics if usertype creation returns an error
func MustUserType(ut *UserType, err error) *UserType {
	if err != nil {
		panic(err)
	}
	return ut
}

// UserType defines types of users
type UserType struct {
	Name      string
	Desc      string
	Szenarios []Szenario
}

// ByUser retruns a Szenario slice defined by the users type
func (c Config) ByUser(u User) ([]Szenario, error) {
	ut, ok := c.userTypes[u.Type()]
	if !ok {
		return nil, fmt.Errorf("No such usertype: %s", u.Type())

	}
	szs := make([]Szenario, len(ut.Szenarios))
	for i, s := range ut.Szenarios {
		s.SetUser(u)
		szs[i] = s
	}
	if len(szs) < 1 {
		slog.Warn("No szenario found for usertype", "user_type", u.Type())
	}
	return szs, nil
}

// ByName return the szario by name, names are case insensitive and if no exact match is found a prefix match is done
func (c Config) ByName(name string) Szenario {
	all := c.allSz.Szenarios
	name = strings.ToLower(name)
	for _, s := range all {
		szname := strings.ToLower(s.Name())
		if szname == name {
			return s
		}
		szname = strings.ReplaceAll(szname, " ", "")
		if szname == name {
			return s
		}
	}
	for _, s := range all {
		szname := strings.ToLower(s.Name())
		if strings.HasPrefix(szname, name) {
			return s
		}
	}
	return nil
}

// GetUserType returns the usertype
func (c Config) GetUserType(name string) *UserType {
	ut, ok := c.userTypes[name]
	if !ok {
		return nil
	}
	return ut
}

// GetUserTypes returns a list of all user types
func (c Config) GetUserTypes() []string {
	utList := make([]string, 0)
	for ut := range c.userTypes {
		utList = append(utList, ut)
	}
	return utList
}

// CreateUserType creates a user type
func (c Config) CreateUserType(n string, desc string) (*UserType, error) {
	ut, found := c.userTypes[n]
	if found {
		ut.Desc = desc
		c.userTypes[n] = ut
		return ut, fmt.Errorf("User type %s already exists: %s", n, ut.Desc)
	}
	ut = &UserType{
		Name:      n,
		Desc:      desc,
		Szenarios: make([]Szenario, 0),
	}
	c.userTypes[n] = ut
	return ut, nil
}

func (c Config) addUserType(ut *UserType, sz Szenario) error {
	us, ok := c.userTypes[ut.Name]
	if !ok {
		return fmt.Errorf("User type %s does not exists", ut)
	}
	us.Szenarios = append(us.Szenarios, sz)
	//userTypes[ut] = us
	return nil
}
