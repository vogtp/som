package user

import (
	"encoding/json"
	"fmt"

	"github.com/suborbital/grav/grav"
	"github.com/vogtp/som/pkg/core"
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
	hcl := core.Get().HCL().Named("user.client")
	hcl.Debugf("Requesting user: %v", name)
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	user := new(User)
	err := p.Send(grav.NewMsg(msgtype.UserRequest, []byte(name))).WaitUntil(defaultTimeout, func(m grav.Message) error {
		hcl.Debugf("Reply for user %s: %T %+v", name, m, m.Data())
		switch m.Type() {
		case msgtype.UserResponse:
			return m.UnmarshalData(user)
		case msgtype.UserError:
			return fmt.Errorf("user backend error: %v", string(m.Data()))
		default:
			return fmt.Errorf("unknown message type %s : %+v", m.Type(), m)
		}

	})
	if err != nil {
		hcl.Errorf("Failed to get user %s: %v", name, err)
		if u, ok := backend.data[name]; ok {
			hcl.Errorf("using local user: %v", name)
			return &u, nil
		}
		return nil, err
	}
	hcl.Debugf("Received user: %v", name)
	return user, nil
}

// Add a user to the store
func (us *client) Add(u *User) error {
	hcl := core.Get().HCL().Named("user.client")
	if err := u.IsValid(); err != nil {
		return fmt.Errorf("user is not valid: %w", err)
	}
	hcl.Debugf("Adding user: %v", u.Name())
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	b, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("cannot marshal user: %w", err)
	}
	msg := grav.NewMsg(msgtype.UserAdd, b)
	var retErr error
	p.Send(msg).WaitUntil(defaultTimeout, func(m grav.Message) error {
		if len(m.Data()) > 0 {
			retErr = fmt.Errorf("server side error: %v", string(m.Data()))
		}
		return nil
	})
	return retErr
}

// List returns a list of all users
func (us *client) List() ([]User, error) {
	hcl := core.Get().HCL().Named("user.client")
	hcl.Debugf("Requesting user list")
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	users := make([]User, 0)
	err := p.Send(grav.NewMsg(msgtype.UserList, nil)).WaitUntil(defaultTimeout, func(m grav.Message) error {
		hcl.Debugf("Reply for userlist: %T %+v", m, m.Data())
		switch m.Type() {
		case msgtype.UserResponse:
			return m.UnmarshalData(&users)
		case msgtype.UserError:
			return fmt.Errorf("user backend error: %v", string(m.Data()))
		default:
			return fmt.Errorf("unknown message type %s : %+v", m.Type(), m)
		}

	})
	if err != nil {
		hcl.Errorf("Failed to get userlist: %v", err)
		return nil, err
	}
	hcl.Debugf("Received users: %v", users)
	return users, nil
}
