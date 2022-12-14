package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/suborbital/grav/grav"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/msgtype"
)

var (
	backend = createBackend()
)

// store stores users and their passwords
type store struct {
	hcl        hcl.Logger
	handlerPod *grav.Pod
	mu         sync.RWMutex
	data       map[string]User
}

// IntialiseStore does the setup for the user store
// starts a goroutine and handles user request in the background
func IntialiseStore() {
	backend.setup()
	backend.start()
	hcl.Warnf("Userstore backend started (key length: %v)", len(core.Keystore.Key()))
}

// createBackend creates a new UserStore
func createBackend() *store {
	return &store{
		hcl:  hcl.New(hcl.WithName("user.store.backend")),
		data: make(map[string]User),
	}
}

func (us *store) setup() {
	c := core.Get()
	us.hcl = c.HCL().Named("user.store.backend")
	if err := us.load(); err != nil {
		us.hcl.Errorf("Cannot load users: %v", err)
	}
}

func (us *store) start() {
	us.handlerPod = core.Get().Bus().Connect()
	us.handlerPod.On(func(m grav.Message) error {
		us.hcl.Tracef("user backend got message: %s %s ID: %v", m.Type(), string(m.Data()), m.UUID())
		switch m.Type() {
		case msgtype.UserRequest:
			return us.getUser(m)
		case msgtype.UserList:
			return us.getUserList(m)
		case msgtype.UserAdd:
			return us.addUser(m)
		case msgtype.UserDelete:
			return us.deleteUser(m)
		case msgtype.UserError:
			return nil
		case msgtype.UserResponse:
			return nil
		default:
			if strings.HasPrefix(m.Type(), "user") {
				us.hcl.Warnf("unhandled user message type: %s -> %v", m.Type(), string(m.Data()))
			}
			return nil
		}
	})
	us.hcl.Tracef("Userstore pod for msg handling: %+v", us.handlerPod)
}

func (us *store) addUser(m grav.Message) error {
	us.hcl.Debug("Requested to add a user")
	_, err := us.storeUserFromMsg(m)
	var s string
	if err != nil {
		us.hcl.Warnf("adding user: %v", err)
		s = err.Error()
	}
	msg := grav.NewMsg(msgtype.UserResponse, []byte(s))
	msg.SetReplyTo(m.UUID())
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	p.Send(msg)

	return err
}

func (us *store) deleteUser(m grav.Message) error {
	name := string(m.Data())
	us.hcl.Warnf("Deleting user %s from store", name)

	var msgTxt string
	msgType := msgtype.UserError
	if _, ok := us.data[name]; ok {
		us.mu.Lock()
		delete(us.data, name)
		us.mu.Unlock()
		if err := us.save(); err != nil {
			us.hcl.Warnf("Cannot save store to delete user %v: %v", name, err)
			msgTxt = fmt.Sprintf("Cannot save store to delete user %v: %v", name, err)
		} else {
			msgTxt = fmt.Sprintf("Deleted %s", name)
			msgType = msgtype.UserResponse
		}
	} else {
		msgTxt = fmt.Sprintf("No such user %s", name)
	}

	msg := grav.NewMsg(msgType, []byte(msgTxt))
	msg.SetReplyTo(m.UUID())
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	p.Send(msg)
	return nil
}

func (us *store) storeUserFromMsg(m grav.Message) (*User, error) {
	u := &User{}
	if err := json.Unmarshal(m.Data(), u); err != nil {
		return nil, fmt.Errorf("adding user: %v", err)
	}
	if err := u.IsValid(); err != nil {
		return nil, fmt.Errorf("new user %v is not valid: %v", u.Name(), err)
	}
	us.mu.Lock()
	if oldUser, ok := us.data[u.Name()]; ok {
		for _, oldPw := range oldUser.History {
			found := false
			for _, newPw := range u.History {
				if string(oldPw.Passwd) == string(newPw.Passwd) {
					found = true
					break
				}
			}
			if found {
				continue
			}
			u.History = append(u.History, oldPw)
		}

	}
	us.data[u.Name()] = *u
	us.mu.Unlock()
	us.hcl.Infof("Added user %s to store", u.Name())
	return u, us.save()
}

func (us *store) getUser(m grav.Message) error {
	name := string(m.Data())
	us.hcl.Debugf("Looking up user %s in store", name)

	msg, err := us.buildUserMsg(name)

	msg.SetReplyTo(m.UUID())
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	p.Send(msg)
	return err
}

func (us *store) buildUserMsg(name string) (grav.Message, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	if u, ok := us.data[name]; ok {
		b, err := json.Marshal(u)
		if err != nil {
			err = fmt.Errorf("cannot marshall user %s: %v", name, err)
			us.hcl.Errorf(err.Error())
			return grav.NewMsg(msgtype.UserError, []byte(err.Error())), err
		}
		return grav.NewMsg(msgtype.UserResponse, b), nil
	}
	return grav.NewMsg(msgtype.UserError, []byte("No such user")), errors.New("no such user")
}

func (us *store) getUserList(m grav.Message) error {
	msg, err := us.buildUserlistMsg()

	msg.SetReplyTo(m.UUID())
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	p.Send(msg)
	return err
}

func (us *store) buildUserlistMsg() (grav.Message, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	users := make([]User, 0, len(us.data))
	for _, u := range us.data {
		users = append(users, u)
	}
	b, err := json.Marshal(users)
	if err != nil {
		err = fmt.Errorf("cannot marshall userlist: %v", err)
		us.hcl.Errorf(err.Error())
		return grav.NewMsg(msgtype.UserError, []byte(err.Error())), err
	}
	return grav.NewMsg(msgtype.UserResponse, b), nil
}

// Get returns the requested user or nil
func (us *store) Get(name string) *User {
	us.mu.RLock()
	defer us.mu.RUnlock()
	u := us.data[name]
	return &u
}

// Add adds a user and encrypts its password
func (us *store) Add(u User, password string) {
	pw := encrypt([]byte(password), core.Keystore.Key())
	us.AddRaw(u, pw)
	hcl.Errorf("Replace with AddRaw(%#v, %#v)\n", u, pw)
}

// AddRaw adds a user with its already encrypted password
func (us *store) AddRaw(u User, password []byte) {
	if len(u.Name()) < 1 {
		us.hcl.Warnf("User must have a name: %v", u)
		return
	}
	defer func() {
		if err := us.save(); err != nil {
			backend.hcl.Warnf("cannot save user store: %v", err)
		}
	}()
	us.mu.Lock()
	defer us.mu.Unlock()
	u.History = []*PwEntry{
		{
			Passwd:  password,
			Created: time.Now(),
		},
	}
	us.data[u.Name()] = u
}
