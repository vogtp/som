package user

import (
	"encoding/json"
	"fmt"

	"github.com/suborbital/grav/grav"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msgtype"
)

// Access to the user store
type Access interface {
	Get(name string) *User
}

var (
	// Store is the userstore client
	Store = createClient()
)

type client struct {
	timeout grav.TimeoutFunc
}

// createClient creates the access level to the user store
func createClient() *client {
	return &client{
		timeout: grav.Timeout(15),
	}
}

// Timeout sets the timeout of the bus in seconds
func (us *client) Timeout(s int) {
	us.timeout = grav.Timeout(s)
}

// Get returns the requested user or nil
func (us *client) Get(name string) (*User, error) {
	slog := core.Get().Log().With(log.Component, "user.client")
	slog.Debug("Requesting user", log.User, name)
	user := new(User)

	d, err := core.Get().AmqpBus().Ask(msgtype.UserRequest, []byte(name))
	if err != nil {
		slog.Warn("Failed to get user", log.User, name, log.Error, err)
		if u, ok := backend.data[name]; ok {
			slog.Error("using local user", log.User, name)
			return &u, nil
		}
		return nil, err
	}
	err = json.Unmarshal(d.Body, user)
	if err != nil {
		return nil, fmt.Errorf("unmarshal user %s response: %w", name, err)
	}

	slog.Debug("Received user", log.User, name)
	return user, nil
}

// Save a user to the store
func (us *client) Save(u *User) error {
	slog := core.Get().Log().With(log.Component, "user.client", log.User, u.Name())
	if err := u.IsValid(); err != nil {
		return fmt.Errorf("user is not valid: %w", err)
	}
	slog.Debug("Saving user")
	b, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("cannot marshal user: %w", err)
	}
	q, err := core.Get().AmqpBus().Ask(msgtype.UserAdd, b)
	if err != nil {
		return fmt.Errorf("save user via bus: %w", err)
	}
	if len(q.Body) > 0 {
		return fmt.Errorf("server side error: %v", string(q.Body))
	}
	return nil
}

// List returns a list of all users
func (us *client) List() ([]User, error) {
	slog := core.Get().Log().With(log.Component, "user.client")
	slog.Debug("Requesting user list")
	users := make([]User, 0)

	d, err := core.Get().AmqpBus().Ask(msgtype.UserList, nil)
	if err != nil {
		slog.Error("Failed to get userlist", log.Error, err)
		return nil, err
	}
	if err := json.Unmarshal(d.Body, &users); err != nil {
		return nil, fmt.Errorf("unmarshal user list: %w", err)
	}
	slog.Debug("Received users", "users", users)
	return users, nil
}

// Delete the user
func (us *client) Delete(name string) (string, error) {
	slog := core.Get().Log().With(log.Component, "user.client", log.User, name)
	slog.Debug("Deleting user")
	d, err := core.Get().AmqpBus().Ask(msgtype.UserDelete, []byte(name))
	msg := string(d.Body)
	return msg, err
}
