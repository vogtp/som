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
	Store          = createClient()
	defaultTimeout = grav.Timeout(15)
)

type client struct {
}

// createClient creates the access level to the user store
func createClient() *client {

	return &client{}
}

// Get returns the requested user or nil
func (us *client) Get(name string) (*User, error) {
	hcl := core.Get().HCL().With(log.Component, "user.client")
	hcl.Debug("Requesting user", log.User, name)
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	user := new(User)
	err := p.Send(grav.NewMsg(msgtype.UserRequest, []byte(name))).WaitUntil(defaultTimeout, func(m grav.Message) error {
		//hcl.Tracef("Reply for user %s: %T %+v", name, m, string(m.Data()))
		switch m.Type() {
		case msgtype.UserResponse:
			return json.Unmarshal(m.Data(), user)
		case msgtype.UserError:
			return fmt.Errorf("user backend error: %v", string(m.Data()))
		default:
			return fmt.Errorf("unknown message type %s : %+v", m.Type(), m)
		}

	})
	if err != nil {
		hcl.Warn("Failed to get user", log.User, name, log.Error, err)
		if u, ok := backend.data[name]; ok {
			hcl.Error("using local user", log.User, name)
			return &u, nil
		}
		return nil, err
	}
	hcl.Debug("Received user", log.User, name)
	return user, nil
}

// Save a user to the store
func (us *client) Save(u *User) error {
	hcl := core.Get().HCL().With(log.Component, "user.client", log.User, u.Name())
	if err := u.IsValid(); err != nil {
		return fmt.Errorf("user is not valid: %w", err)
	}
	hcl.Debug("Saving user")
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	b, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("cannot marshal user: %w", err)
	}
	msg := grav.NewMsg(msgtype.UserAdd, b)
	var retErr error
	err = p.Send(msg).WaitUntil(defaultTimeout, func(m grav.Message) error {
		if len(m.Data()) > 0 {
			retErr = fmt.Errorf("server side error: %v", string(m.Data()))
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("cannot send user msg: %w", err)
	}
	return retErr
}

// List returns a list of all users
func (us *client) List() ([]User, error) {
	hcl := core.Get().HCL().With(log.Component, "user.client")
	hcl.Debug("Requesting user list")
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	users := make([]User, 0)
	err := p.Send(grav.NewMsg(msgtype.UserList, nil)).WaitUntil(defaultTimeout, func(m grav.Message) error {
		//hcl.Tracef("Reply for userlist: %T %+v", m, string(m.Data()))
		switch m.Type() {
		case msgtype.UserResponse:
			return json.Unmarshal(m.Data(), &users)
		case msgtype.UserError:
			return fmt.Errorf("user backend error: %v", string(m.Data()))
		default:
			return fmt.Errorf("unknown message type %s : %+v", m.Type(), m)
		}

	})
	if err != nil {
		hcl.Error("Failed to get userlist", log.Error, err)
		return nil, err
	}
	hcl.Debug("Received users", "users", users)
	return users, nil
}

// Delete the user
func (us *client) Delete(name string) (string, error) {
	hcl := core.Get().HCL().With(log.Component, "user.client", log.User, name)
	hcl.Debug("Deleting user")
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	msg := ""
	err := p.Send(grav.NewMsg(msgtype.UserDelete, []byte(name))).WaitUntil(defaultTimeout, func(m grav.Message) error {
		//hcl.Tracef("Reply for user %s: %T %+v", name, m, string(m.Data()))
		switch m.Type() {
		case msgtype.UserResponse:
			msg = string(m.Data())
			return nil
		case msgtype.UserError:
			return fmt.Errorf("user backend error: %v", string(m.Data()))
		default:
			return fmt.Errorf("unknown message type %s : %+v", m.Type(), m)
		}

	})

	return msg, err
}
