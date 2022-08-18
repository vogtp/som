package user

import (
	"testing"

	"github.com/vogtp/som/pkg/core"
)

var (
	testKey = []byte("mySuperSecretKey")
	pass    = "pass"
	encPass = []byte{10, 157, 71, 137, 68, 85, 19, 75, 205, 79, 237, 208, 191, 75, 96, 89, 129, 111, 111, 9, 214, 74, 176, 34, 177, 185, 179, 93, 4, 148, 243, 84}
)

func TestUserStorage_Get(t *testing.T) {
	store := createBackend()
	core.Keystore.Add(testKey)
	store.Add(User{Username: "user1", Mail: "mail1", UserType: "stud"}, pass)
	store.AddRaw(User{Username: "user2", Mail: "mail2", UserType: "stud"}, encPass)
	u1 := store.Get("user1")
	u2 := store.Get("user2")
	if u1.Username != "user1" {
		t.Error("Wrong name")
	}
	if u1.Password() != u2.Password() {
		t.Error("Password messup")
	}
}
