package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vogtp/som/pkg/core"
)

// User stores a user and its encrypted password
type User struct {
	Username string `json:"name"`
	Mail     string `json:"email"`
	Longname string `json:"displayname"`
	Passwd   []byte `json:"payload"`
	UserType string `json:"type"`
}

// Name returns the name
func (u *User) Name() string {
	if u == nil {
		return "unknown user"
	}
	return u.Username
}

// Email returns the email address
func (u *User) Email() string {
	return u.Mail
}

// DisplayName returns the display name
func (u *User) DisplayName() string {
	return u.Longname
}

// Type returns what kind of user it is
func (u *User) Type() string {
	return u.UserType
}

// Password decrypts the password
func (u *User) Password() string {
	return string(decrypt(u.Passwd, core.Keystore.Key()))
}

// SetPassword encrypts the password
func (u *User) SetPassword(pw string) {
	u.Passwd = encrypt([]byte(pw), core.Keystore.Key())
}

// String implements stringer
func (u User) String() string {
	return fmt.Sprintf("%-30s %-30s %-10s", u.Name(), u.Email(), u.Type())
}

// IsValid checks if all needed fields are set
func (u User) IsValid() error {
	if len(u.Username) < 1 {
		return errors.New("users must have a name")
	}
	if len(u.Mail) < 1 {
		return errors.New("users must have a email")
	}

	if !strings.Contains(u.Mail, "@") {
		return fmt.Errorf("user email %q must be valid", u.Mail)
	}
	if len(u.UserType) < 1 {
		return errors.New("users must have a type")
	}
	if len(u.Passwd) < 1 {
		return errors.New("users must have a password")
	}
	if len(u.Password()) < 1 {
		return errors.New("users password is not decryptable")
	}
	return nil
}
