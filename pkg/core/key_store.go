package core

import "github.com/vogtp/go-hcl"

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
	if len(ks.data) < 10 {
		hcl := hcl.Named("som.keystore")
		hcl.Errorf("Keystore is not initialised! Key len: %v", len(ks.data))
		hcl.Error("Use core.Keystore.Add(...)")
		panic("No store key")
	}
	return ks.data
}
