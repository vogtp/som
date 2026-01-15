package user

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/log"
)

const (
	dbFile = "userstore.db"
	dbPerm = 0600
)

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
		us.log.Error("Could not open cipher reader to load file", log.Error, err, "db_name", dbFile)
		panic(err)
	}

	if err = gob.NewDecoder(r).Decode(&us.data); err != nil {
		return fmt.Errorf("cannot decode users from gob: %v", err)
	}
	us.log.Info("Loaded users from", "count", len(us.data), "file", dbFile)
	if err := us.mirgrate(); err != nil {
		us.log.Warn("Cannot mirgrate", "file", dbFile, log.Error, err)
	}
	return nil
}

func (us *store) save() error {
	us.mu.Lock()
	defer us.mu.Unlock()
	if err := us.backup(); err != nil {
		return fmt.Errorf("create userstore backup: %w", err)
	}
	us.cleanupPasswords()
	f, err := os.OpenFile(dbFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, dbPerm)
	if err != nil {
		us.log.Error("cannot open gob file", "file", dbFile, log.Error, err)
		return fmt.Errorf("cannot open gob file %s: %v", dbFile, err)
	}
	defer f.Close()
	w, err := EncryptedWriter(string(core.Keystore.Key()), f)
	if err != nil {
		us.log.Error("Could not open cipher reader to load file", log.Error, err, "db_name", dbFile)
		panic(err)
	}
	if err = gob.NewEncoder(w).Encode(&us.data); err != nil {
		us.log.Error("cannot encode users to gob", "file", dbFile, log.Error, err)
		return fmt.Errorf("cannot encode users to gob: %v", err)
	}
	us.log.Info("Saved users", "count", len(us.data), "file", dbFile)
	return nil
}

func (us *store) cleanupPasswords() {
	for n, u := range us.data {
		u := u
		u.deleteOldPasswords()
		us.data[n] = u
	}
}

func (us *store) backup() error {
	f, err := os.Open(dbFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	b, err := os.OpenFile(fmt.Sprintf("%s.bck", dbFile), os.O_CREATE|os.O_RDWR|os.O_TRUNC, dbPerm)
	if err != nil {
		return err
	}
	defer b.Close()
	bi, err := b.Stat()
	if err != nil {
		return err
	}
	if fi.Size() < 1 && bi.Size() > 0 {
		return fmt.Errorf("no db %s (size: %v) but a backup %s (size: %v)", dbFile, fi.Size(), bi.Name(), bi.Size())
	}
	_, err = io.Copy(b, f)
	return err
}
