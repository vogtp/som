package core

import (
	"os"

	"github.com/vogtp/som/pkg/core/log"
)

var (
	// Keystore stores the keys
	Keystore = &keyStore{}
)

// keyStore holds the keys
type keyStore struct {
	data []byte
}

// Add a key to store
func (ks *keyStore) Add(k []byte) {
	ks.data = k
}

// Key retruns the key of the store
func (ks *keyStore) Key() []byte {
	if len(ks.data) < 1 {
		log := log.New("som.keystore")
		log.Error("Keystore is not initialised!", "key_len", len(ks.data))
		log.Error("Use core.Keystore.Add(...)")
		os.Exit(1)
	}
	return ks.data
}
