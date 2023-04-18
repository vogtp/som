package user

import (
	"time"
)

func (us *store) mirgrate() error {
	for name, user := range us.data {
		if len(user.History) < 1 {
			us.log.Info("Mirgating PW history", "name", name)
			user.History = make([]*PwEntry, 1)
			user.History[0] = &PwEntry{Passwd: user.DeprecatedPasswd, Created: time.Now()}
			us.data[name] = user
		}
	}
	us.mu.Unlock()
	err := us.save()
	us.mu.Lock()
	return err
}
