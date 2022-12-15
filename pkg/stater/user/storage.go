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
		us.hcl.Errorf("EncryptedReader: %v", err)
		panic(err)
	}

	if err = gob.NewDecoder(r).Decode(&us.data); err != nil {
		return fmt.Errorf("cannot decode users from gob: %v", err)
	}
	us.hcl.Infof("Loaded %d users from %s", len(us.data), dbFile)
	if err := us.mirgrate(); err != nil {
		us.hcl.Warnf("Cannot mirgrate %s: %v", dbFile, err)
	}
	return nil
}

func (us *store) save() error {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.cleanupPasswords()
	f, err := os.OpenFile(dbFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		us.hcl.Errorf("cannot open gob file %s: %v", dbFile, err)
		return fmt.Errorf("cannot open gob file %s: %v", dbFile, err)
	}
	defer f.Close()
	w, err := EncryptedWriter(string(core.Keystore.Key()), f)
	if err != nil {
		us.hcl.Errorf("EncryptedWriter: %v", err)
		panic(err)
	}
	if err = gob.NewEncoder(w).Encode(&us.data); err != nil {
		us.hcl.Errorf("cannot encode users to gob: %v", err)
		return fmt.Errorf("cannot encode users to gob: %v", err)
	}
	us.hcl.Infof("Saved %d users to %s", len(us.data), dbFile)
	return nil
}

func (us *store) cleanupPasswords() {
	for n, u := range us.data {
		u.deleteOldPasswords()
		us.data[n] = u
	}
}
