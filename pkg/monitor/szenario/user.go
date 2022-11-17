package szenario

import "time"

// User stores a user and its encrypted password
type User interface {
	// Name returns the name
	Name() string
	// Email returns the email address
	Email() string

	// DisplayName returns the display name
	DisplayName() string

	// Type returns what kind of user it is
	Type() string

	// Password decrypts the password
	Password() string

	// NextPassword increases the password index and returns the decrypted PW
	// retruns empty string "" if no more passwords are present
	NextPassword() string

	// NewPassword generates a new password
	// it does not store the password
	NewPassword() (string, error)

	// PasswordHistoryCount returns the number of PW in the history
	PasswordHistoryCount() int

	// PasswordCreated returns the time when the password was created
	PasswordCreated() time.Time

	// PasswordLastUse returns the time when the password was last accessed
	PasswordLastUse() time.Time

	// SetPassword encrypts the password
	SetPassword(pw string)

	// String implements stringer
	String() string

	// IsValid checks if all needed fields are set
	IsValid() error
}
