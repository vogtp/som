package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/suborbital/grav/grav"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/log"
)

var (
	backend = createBackend()
)

// store stores users and their passwords
type store struct {
	log        *slog.Logger
	handlerPod *grav.Pod
	mu         sync.RWMutex
	data       map[string]User
}

// IntialiseStore does the setup for the user store
// starts a goroutine and handles user request in the background
func IntialiseStore(ctx context.Context) {
	backend.setup()
	backend.start(ctx)
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

func (us *store) start(ctx context.Context) {
	// routingKey := "som.user.#"
	// err := core.Get().AmqpBus().Answer(ctx, routingKey, func(routingKey string, d amqp.Delivery) ([]byte, error) {
	// 	us.log.Debug("user backend got message", "type", routingKey, "data", string(d.Body))
	// 	switch routingKey {
	// 	case msgtype.UserRequest:
	// 		return us.getUser(d)
	// 	case msgtype.UserList:
	// 		return us.getUserList(d)
	// 	case msgtype.UserAdd:
	// 		return us.addUser(d)
	// 	case msgtype.UserDelete:
	// 		return us.deleteUser(d)
	// 	case msgtype.UserError:
	// 		return nil, nil
	// 	default:
	// 		if strings.HasPrefix(routingKey, "user") {
	// 			us.log.Warn("unhandled user message type", "type", routingKey, "data", string(d.Body))
	// 		}
	// 		return nil, nil
	// 	}
	// })
	// if err != nil {
	// 	us.log.Error("Cannot listen on bus", "routingKey", routingKey, log.Error, err)
	// }
	us.log.Debug("Userstore pod for msg handling", "pod", us.handlerPod)
}

func (us *store) addUser(d amqp.Delivery) ([]byte, error) {
	us.log.Debug("Requested to add a user")
	_, err := us.storeUserFromMsg(d)
	var s string
	if err != nil {
		us.log.Warn("adding user", log.Error, err)
		s = err.Error()
	}
	return []byte(s), err
}

func (us *store) deleteUser(d amqp.Delivery) ([]byte, error) {
	name := string(d.Body)
	us.log.Warn("Deleting user from store", log.User, name)

	var msgTxt string
	if _, ok := us.data[name]; ok {
		us.mu.Lock()
		delete(us.data, name)
		us.mu.Unlock()
		if err := us.save(); err != nil {
			us.log.Warn("Cannot save store to delete user", log.User, name, log.Error, err)
			msgTxt = fmt.Sprintf("Cannot save store to delete user %v: %v", name, err)
		} else {
			msgTxt = fmt.Sprintf("Deleted %s", name)
		}
	} else {
		msgTxt = fmt.Sprintf("No such user %s", name)
	}
	return []byte(msgTxt), nil
}

func (us *store) storeUserFromMsg(d amqp.Delivery) (*User, error) {
	u := &User{}
	if err := json.Unmarshal(d.Body, u); err != nil {
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

func (us *store) getUser(d amqp.Delivery) ([]byte, error) {
	name := string(d.Body)
	us.log.Debug("Looking up user in store", log.User, name)

	return us.buildUserMsg(name)
}

func (us *store) buildUserMsg(name string) ([]byte, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	if u, ok := us.data[name]; ok {
		b, err := json.Marshal(u)
		if err != nil {
			err = fmt.Errorf("cannot marshall user %s: %v", name, err)
			us.log.Error("cannot marshall user", log.Error, err.Error(), log.User, name)
			return []byte(err.Error()), err
		}
		return b, nil
	}
	return []byte("No such user"), errors.New("no such user")
}

func (us *store) getUserList(d amqp.Delivery) ([]byte, error) {
	return us.buildUserlistMsg()
}

func (us *store) buildUserlistMsg() ([]byte, error) {
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
		return []byte(err.Error()), err
	}
	return b, nil
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
