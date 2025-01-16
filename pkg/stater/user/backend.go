package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"log/slog"

	"github.com/suborbital/e2core/foundation/bus/bus"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msgtype"
)

var (
	backend = createBackend()
)

// store stores users and their passwords
type store struct {
	log        *slog.Logger
	handlerPod *bus.Pod
	mu         sync.RWMutex
	data       map[string]User
}

// IntialiseStore does the setup for the user store
// starts a goroutine and handles user request in the background
func IntialiseStore() {
	backend.setup()
	backend.start()
	slog.Warn("Userstore backend started", "key_len", len(core.Keystore.Key()))
}

// createBackend creates a new UserStore
func createBackend() *store {
	return &store{
		log:  log.New("user.store.backend"),
		data: make(map[string]User),
	}
}

func (us *store) setup() {
	c := core.Get()
	us.log = c.Log().With(log.Component, "user.store.backend")
	if err := us.load(); err != nil {
		us.log.Error("Cannot load users", log.Error, err)
	}
}

func (us *store) start() {
	us.handlerPod = core.Get().Bus().Connect()
	us.handlerPod.On(func(m bus.Message) error {
		us.log.Debug("user backend got message", "type", m.Type(), "data", string(m.Data()), "uuid", m.UUID())
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
				us.log.Warn("unhandled user message type", "type", m.Type(), "data", string(m.Data()))
			}
			return nil
		}
	})
	us.log.Debug("Userstore pod for msg handling", "pod", us.handlerPod)
}

func (us *store) addUser(m bus.Message) error {
	us.log.Debug("Requested to add a user")
	_, err := us.storeUserFromMsg(m)
	var s string
	if err != nil {
		us.log.Warn("adding user", log.Error, err)
		s = err.Error()
	}
	msg := bus.NewMsg(msgtype.UserResponse, []byte(s))
	msg.SetReplyTo(m.UUID())
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	p.Send(msg)

	return err
}

func (us *store) deleteUser(m bus.Message) error {
	name := string(m.Data())
	us.log.Warn("Deleting user from store", log.User, name)

	var msgTxt string
	msgType := msgtype.UserError
	if _, ok := us.data[name]; ok {
		us.mu.Lock()
		delete(us.data, name)
		us.mu.Unlock()
		if err := us.save(); err != nil {
			us.log.Warn("Cannot save store to delete user", log.User, name, log.Error, err)
			msgTxt = fmt.Sprintf("Cannot save store to delete user %v: %v", name, err)
		} else {
			msgTxt = fmt.Sprintf("Deleted %s", name)
			msgType = msgtype.UserResponse
		}
	} else {
		msgTxt = fmt.Sprintf("No such user %s", name)
	}

	msg := bus.NewMsg(msgType, []byte(msgTxt))
	msg.SetReplyTo(m.UUID())
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	p.Send(msg)
	return nil
}

func (us *store) storeUserFromMsg(m bus.Message) (*User, error) {
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
	us.log.Info("Added user to store", log.User, u.Name())
	return u, us.save()
}

func (us *store) getUser(m bus.Message) error {
	name := string(m.Data())
	us.log.Debug("Looking up user in store", log.User, name)

	msg, err := us.buildUserMsg(name)

	msg.SetReplyTo(m.UUID())
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	p.Send(msg)
	return err
}

func (us *store) buildUserMsg(name string) (bus.Message, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	if u, ok := us.data[name]; ok {
		b, err := json.Marshal(u)
		if err != nil {
			err = fmt.Errorf("cannot marshall user %s: %v", name, err)
			us.log.Error("cannot marshall user", log.Error, err.Error(), log.User, name)
			return bus.NewMsg(msgtype.UserError, []byte(err.Error())), err
		}
		return bus.NewMsg(msgtype.UserResponse, b), nil
	}
	return bus.NewMsg(msgtype.UserError, []byte("No such user")), errors.New("no such user")
}

func (us *store) getUserList(m bus.Message) error {
	msg, err := us.buildUserlistMsg()

	msg.SetReplyTo(m.UUID())
	p := core.Get().Bus().Connect()
	defer p.Disconnect()
	p.Send(msg)
	return err
}

func (us *store) buildUserlistMsg() (bus.Message, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	users := make([]User, 0, len(us.data))
	for _, u := range us.data {
		users = append(users, u)
	}
	b, err := json.Marshal(users)
	if err != nil {
		err = fmt.Errorf("cannot marshall userlist: %v", err)
		us.log.Error("Cannot marshall user list", log.Error, err.Error())
		return bus.NewMsg(msgtype.UserError, []byte(err.Error())), err
	}
	return bus.NewMsg(msgtype.UserResponse, b), nil
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
	fmt.Printf("Replace with AddRaw(%#v, %#v)\n", u, pw)
}

// AddRaw adds a user with its already encrypted password
func (us *store) AddRaw(u User, password []byte) {
	if len(u.Name()) < 1 {
		us.log.Warn("User must have a name", log.User, u)
		return
	}
	defer func() {
		if err := us.save(); err != nil {
			backend.log.Warn("cannot save user store", log.Error, err)
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
