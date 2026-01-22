package user_test

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/stater"
	"github.com/vogtp/som/pkg/stater/user"
)

var (
	testKey = []byte("mySuperSecretKey")
	pass    = "pass"
	encPass = []byte{10, 157, 71, 137, 68, 85, 19, 75, 205, 79, 237, 208, 191, 75, 96, 89, 129, 111, 111, 9, 214, 74, 176, 34, 177, 185, 179, 93, 4, 148, 243, 84}
	u1      = &user.User{Username: "user1", Mail: "mail1@test.net", UserType: "stud"}
	u2      = &user.User{Username: "user2", Mail: "mail2@test.net", UserType: "stud"}
)

func testInit(t *testing.T) func() {
	core.Keystore.Add(testKey)
	close, err := stater.Run(t.Context(), "user.test")
	if err != nil {
		t.Fatalf("Cannot start stater: %v", err)
	}
	log.Level.Set(slog.LevelDebug)
	user.Store.Timeout(1)
	u1.SetPassword(pass)
	u2.SetPassword(pass)
	if err := user.Store.Save(u1); err != nil {
		t.Fatalf("cannot save u1: %v", err)
	}
	if err := user.Store.Save(u2); err != nil {
		t.Fatalf("cannot save u2: %v", err)
	}
	return close
}

func TestUser_List(t *testing.T) {
	close := testInit(t)
	defer close()

	users, err := user.Store.List()
	if err != nil {
		t.Fatalf("Cannot get users: %v", err)
	}
	if len(users) < 2 {
		t.Fatalf("Did not find users: %v", len(users))
	}

}

func TestUser_Get(t *testing.T) {
	close := testInit(t)
	defer close()

	checkUser(t, u1)
	checkUser(t, u2)
	if u, err := user.Store.Get("not.existing.user"); err == nil {
		t.Fatalf("found non existing user: %v", u)
	}
}

func TestUser_Del(t *testing.T) {
	close := testInit(t)
	defer close()

	checkUser(t, u1)
	checkUser(t, u2)
	if _, err := user.Store.Delete(u1.Name()); err != nil {
		t.Errorf("Cannot delete user: %v", err)
	}
	if u, err := user.Store.Get(u1.Name()); err == nil {
		t.Fatalf("found deleted user: %v", u)
	}
	if _, err := user.Store.Delete(u2.Name()); err != nil {
		t.Errorf("Cannot delete user: %v", err)
	}
	if u, err := user.Store.Get(u2.Name()); err == nil {
		t.Fatalf("found deleted user: %v", u)
	}
	users, err := user.Store.List()
	if err != nil {
		t.Fatalf("Cannot get users: %v", err)
	}
	if len(users) > 0 {
		t.Fatalf("Found users all should have been deleted: %v", len(users))
	}
}

func checkUser(t *testing.T, su *user.User) {
	u, err := user.Store.Get(su.Username)
	if err != nil {
		t.Fatalf("cannot load user: %s", err)
	}
	if reflect.DeepEqual(su, u) {
		t.Fatal("Retrieved user not correct")
	}
}
