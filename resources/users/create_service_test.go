package users_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zoommix/fasthttp_template/resources/users"
	pg "github.com/zoommix/fasthttp_template/store"
)

// TestMain ...
func TestMain(m *testing.M) {
	defer pg.Close()
	code := m.Run()
	os.Exit(code)
}

// TestCreatesIfUserDataIsValid ...
func TestCreatesIfUserDataIsValid(t *testing.T) {
	defer pg.TearDown("users")

	usr := userData()
	u, err := users.CreateUser(usr)
	u = u.Reload()

	assert.NoError(t, err, "No error while creation")
	assert.NotNil(t, u, "User obj is not nil")

	if u != nil {
		assert.NotNil(t, u.ID, "ID is not nil")
		assert.Equal(t, usr.Email, u.Email, "Expected email")
		assert.Equal(t, usr.FirstName, u.FirstName, "Expected first name")
		assert.Equal(t, usr.LastName, u.LastName, "Expected last name")
	}
}

// TestDoesNotCreateUserIfDataInvalid ...
func TestDoesNotCreateUserIfDataInvalid(t *testing.T) {
	defer pg.TearDown("users")

	usr := userData()

	// create user to trigger uniq index error
	users.CreateUser(usr)

	u, err := users.CreateUser(usr)

	assert.Error(t, err, "No error while creation")
	assert.Nil(t, u, "User is not nil")
}

func userData() *users.User {
	return &users.User{
		Email:     "zoommix@ex.ua",
		FirstName: "Roman",
		LastName:  "Huk",
		Password:  "qwertyui",
	}
}

func createUser() *users.User {
	data := userData()

	u, _ := users.CreateUser(data)
	u.PasswordDigest = users.EncodePassword(data.Password)
	u.Password = "" // should be empty to use it as expected result

	return u
}
