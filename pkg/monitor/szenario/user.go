package szenario

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

	// SetPassword encrypts the password
	SetPassword(pw string)

	// String implements stringer
	String() string

	// IsValid checks if all needed fields are set
	IsValid() error
}
