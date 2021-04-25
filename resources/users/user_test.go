package users_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zoommix/fasthttp_template/resources/users"
	pg "github.com/zoommix/fasthttp_template/store"
)

// TestFindUser ...
func TestFindUser(t *testing.T) {
	defer pg.TearDown("users")

	premadeUsr := createUser()

	var testCases = []struct {
		id   int
		err  error
		user *users.User
	}{
		{premadeUsr.ID, nil, premadeUsr},
		{-1, errors.New("user is not found"), nil},
	}

	for _, c := range testCases {
		u, err := users.FindUser(c.id)

		assert.Equal(t, c.err, err, "No error while creation")
		assert.Equal(t, c.user, u, "returns expected fields")
	}
}

// TestFindUserByEmail ...
func TestFindUserByEmail(t *testing.T) {
	defer pg.TearDown("users")

	premadeUsr := createUser()

	var testCases = []struct {
		email string
		err   error
		user  *users.User
	}{
		{premadeUsr.Email, nil, premadeUsr},
		{"some@email.com", errors.New("user is not found"), nil},
	}

	for _, c := range testCases {
		u, err := users.FindUserByEmail(c.email)

		assert.Equal(t, c.err, err, "No error while creation")
		assert.Equal(t, c.user, u, "returns expected fields")
	}
}

// TestEmailExists ...
func TestEmailExists(t *testing.T) {
	defer pg.TearDown("users")

	premadeUsr := createUser()
	testCases := []struct {
		email   string
		isExist bool
	}{
		{premadeUsr.Email, true},
		{"some_random@email.com", false},
	}

	for _, c := range testCases {
		isExist := users.EmailExists(c.email)

		assert.Equal(t, c.isExist, isExist, "Expected boolean value")
	}
}

// TestReload ...
func TestReload(t *testing.T) {
	defer pg.TearDown("users")

	premadeUsr := createUser()

	var testCases = []struct {
		user     *users.User
		expected *users.User
	}{
		{&users.User{ID: premadeUsr.ID}, premadeUsr},
		{&users.User{ID: -1}, nil},
	}

	for _, c := range testCases {
		u := c.user.Reload()

		assert.Equal(t, c.expected, u, "Returns expected user data")
	}
}
