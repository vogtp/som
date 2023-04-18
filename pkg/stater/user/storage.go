package user

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/vogtp/som/pkg/core"
)

const dbFile = "userstore.db"

func (us *store) load() error {
	us.mu.Lock()
	defer us.mu.Unlock()
	f, err := os.OpenFile(dbFile, os.O_RDONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot open gob file %s: %v", dbFile, err)
	}
	defer f.Close()
	r, err := EncryptedReader(string(core.Keystore.Key()), f)
	if err != nil {
		us.log.Error("EncryptedReader", "error", err)
		panic(err)
	}

	if err = gob.NewDecoder(r).Decode(&us.data); err != nil {
		return fmt.Errorf("cannot decode users from gob: %v", err)
	}
	us.log.Info("Loaded users from", "count", len(us.data), "file", dbFile)
	if err := us.mirgrate(); err != nil {
		us.log.Warn("Cannot mirgrate", "file", dbFile, "error", err)
	}
	return nil
}

func (us *store) save() error {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.cleanupPasswords()
	f, err := os.OpenFile(dbFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		us.log.Error("cannot open gob file", "file", dbFile, "error", err)
		return fmt.Errorf("cannot open gob file %s: %v", dbFile, err)
	}
	defer f.Close()
	w, err := EncryptedWriter(string(core.Keystore.Key()), f)
	if err != nil {
		us.log.Error("EncryptedWriter", "file", dbFile, "error", err)
		panic(err)
	}
	if err = gob.NewEncoder(w).Encode(&us.data); err != nil {
		us.log.Error("cannot encode users to gob", "file", dbFile, "error", err)
		return fmt.Errorf("cannot encode users to gob: %v", err)
	}
	us.log.Info("Saved users", "count", len(us.data), "file", dbFile)
	return nil
}

func (us *store) cleanupPasswords() {
	for n, u := range us.data {
		u.deleteOldPasswords()
		us.data[n] = u
	}
}
