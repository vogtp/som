package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vogtp/som/pkg/core"
)

// User stores a user and its encrypted password
type User struct {
	Username string     `json:"name"`
	Mail     string     `json:"email"`
	Longname string     `json:"displayname"`
	Passwd   []byte     `json:"payload"`
	History  []*PwEntry `json:"history"`
	UserType string     `json:"type"`
	pwIdx    int
}

// PwEntry stores pw history
type PwEntry struct {
	Passwd  []byte    `json:"payload"`
	Created time.Time `json:"created"`
	LastUse time.Time `json:"last_use"`
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

// NextPassword increases the password index and returns the decrypted PW
func (u *User) NextPassword() string {
	if !(u.pwIdx+1 < len(u.History)) {
		return ""
	}
	u.pwIdx++
	cur := u.History[u.pwIdx]
	cur.LastUse = time.Now()
	return string(decrypt(cur.Passwd, core.Keystore.Key()))
}

// Password decrypts the password
func (u *User) Password() string {
	// if !(u.pwIdx < len(u.History)) {
	// 	return ""
	// }
	// cur := u.History[u.pwIdx]
	// cur.LastUse = time.Now()
	// return string(decrypt(cur.Passwd, core.Keystore.Key()))
	return string(decrypt(u.Passwd, core.Keystore.Key()))
}

// PasswordCreated returns the time when the password was created
func (u *User) PasswordCreated() time.Time {
	return u.History[u.pwIdx].Created
}

// PasswordLastUse returns the time when the password was last accessed
func (u *User) PasswordLastUse() time.Time {
	return u.History[u.pwIdx].LastUse
}

// SetPassword encrypts the password
func (u *User) SetPassword(pw string) {
	u.Passwd = encrypt([]byte(pw), core.Keystore.Key())
	// pe := PwEntry{
	// 	Passwd:  encrypt([]byte(pw), core.Keystore.Key()),
	// 	Created: time.Now(),
	// }
	// u.History = append([]*PwEntry{&pe}, u.History...)
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
