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
	Username         string     `json:"name"`
	Mail             string     `json:"email"`
	Longname         string     `json:"displayname"`
	DeprecatedPasswd []byte     `json:"payload"` // Deprecated: use History instead. (since v0.11.0)
	History          []*PwEntry `json:"history"`
	UserType         string     `json:"type"`
	pwIdx            int
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
// retruns empty string "" if no more passwords are present
func (u *User) NextPassword() string {
	if !(u.pwIdx+1 < len(u.History)) {
		return ""
	}
	u.pwIdx++
	cur := u.History[u.pwIdx]
	return string(decrypt(cur.Passwd, core.Keystore.Key()))
}

// Password decrypts the password
func (u *User) Password() string {
	if !(u.pwIdx < len(u.History)) {
		return ""
	}
	cur := u.History[u.pwIdx]
	return string(decrypt(cur.Passwd, core.Keystore.Key()))
	// return string(decrypt(u.Passwd, core.Keystore.Key()))
}

// LoginSuccessfull sets the last use of the password
func (u *User) LoginSuccessfull() {
	u.History[u.pwIdx].LastUse = time.Now()
}

// PasswordHistoryCount returns the number of PW in the history
func (u User) PasswordHistoryCount() int {
	return len(u.History)
}

// PasswordCreated returns the time when the password was created
func (u *User) PasswordCreated() time.Time {
	return u.History[u.pwIdx].Created
}

// PasswordLastUse returns the time when the password was last accessed
func (u *User) PasswordLastUse() time.Time {
	return u.History[u.pwIdx].LastUse
}

// ResetPasswordIndex start with the first password
// and reset the number of failed logins to 0
func (u *User) ResetPasswordIndex() {
	u.pwIdx = 0
}

// NumPasswdChg number of times the password was changed
func (u *User) NumPasswdChg(d time.Duration) int {
	cnt := 0
	for _, pw := range u.History {
		if time.Since(pw.Created) < d {
			cnt++
		}
	}
	return cnt
}

// FailedLogins the number of failed logins
func (u User) FailedLogins() int {
	return u.pwIdx
}

func (u *User) deleteOldPasswords() {
	if u.PasswordLastUse().IsZero() {
		return
	}
	hcl := core.Get().HCL().Named("passwordCleanup")
	curAge := time.Since(u.PasswordLastUse())
	if curAge > time.Hour {
		return
	}
	hist := make([]*PwEntry, 0)
	for i := 0; i < len(u.History); i++ {
		if time.Since(u.History[i].LastUse) > 3*24*time.Hour {
			continue
		}
		hist = append(hist, u.History[i])
	}
	if len(hist) < u.pwIdx+1 || len(hist) < 5 || len(hist) >= len(u.History) {
		return
	}
	hcl.Infof("%s deleted old passwords: %v -> %v", u.Name(), len(u.History), len(hist))
	u.History = hist
}

// SetPassword encrypts the password
func (u *User) SetPassword(pw string) {
	// u.Passwd = encrypt([]byte(pw), core.Keystore.Key())
	pe := PwEntry{
		Passwd:  encrypt([]byte(pw), core.Keystore.Key()),
		Created: time.Now(),
	}
	if len(u.History) < 1 || string(u.History[0].Passwd) != string(pe.Passwd) {
		u.History = append([]*PwEntry{&pe}, u.History...)
	}
	core.Get().HCL().Warnf("Change password of user %s", u.Name())
}

// Save the user to the store
func (u *User) Save() error {
	if err := Store.Save(u); err != nil {
		return fmt.Errorf("cannot save %s: %w", u.Name(), err)
	}
	return nil
}

// String implements stringer
func (u User) String() string {
	return fmt.Sprintf("%-30s %-30s %-10s (Password History %2d)", u.Name(), u.Email(), u.Type(), u.PasswordHistoryCount())
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
	if len(u.Password()) < 1 {
		return errors.New("users must have a password")
	}
	return nil
}
