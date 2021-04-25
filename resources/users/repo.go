package users

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	pg "github.com/zoommix/fasthttp_template/store"
	"github.com/zoommix/fasthttp_template/utils"
)

// CreateUser ...
func CreateUser(u *User) (*User, error) {
	passMD5 := EncodePassword(u.Password)

	err := pg.
		DB.
		QueryRow(
			context.Background(),
			"INSERT INTO users (email, first_name, last_name, password_digest, created_at, updated_at)"+
				"VALUES ($1, $2, $3, $4, now() AT TIME ZONE 'UTC', now() AT TIME ZONE 'UTC')"+
				"RETURNING id, extract(epoch from created_at)::integer",
			strings.ToLower(u.Email),
			u.FirstName,
			u.LastName,
			passMD5,
		).
		Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		utils.LogInfo(fmt.Sprintf("Unable to create user with email=%+v", u))

		return nil, errors.New(strings.Replace(err.Error(), `"`, `'`, -1))
	}

	utils.LogInfo(fmt.Sprintf("Created user with ID=%d, email=%s", u.ID, u.Email))

	return u, nil
}

// EncodePassword ...
func EncodePassword(p string) string {
	sha := sha256.Sum256([]byte(p))
	return base64.StdEncoding.EncodeToString(sha[:])
}

// FindUser finds user by ID
func FindUser(ID int) (*User, error) {
	u := &User{}

	if err := pg.DB.QueryRow(
		context.Background(),
		"SELECT id, first_name, last_name, email, password_digest, "+
			"extract(epoch from created_at)::integer created_at "+
			"FROM users "+
			"WHERE users.id = $1 "+
			"LIMIT 1",
		ID,
	).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.PasswordDigest,
		&u.CreatedAt,
	); err != nil && err.Error() != ignoredSQLError {
		utils.LogError(err.Error())
	}

	if u.ID == 0 {
		return nil, errors.New("user is not found")
	}

	return u, nil
}

// FindUserByEmail finds user by ID
func FindUserByEmail(email string) (*User, error) {
	u := &User{}

	if err := pg.DB.QueryRow(
		context.Background(),
		"SELECT id, first_name, last_name, email, password_digest, "+
			"extract(epoch from created_at)::integer created_at "+
			"FROM users "+
			"WHERE users.email = $1 "+
			"LIMIT 1",
		email,
	).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.PasswordDigest,
		&u.CreatedAt,
	); err != nil && err.Error() != ignoredSQLError {
		utils.LogError(err.Error())
	}

	if u.ID == 0 {
		return nil, errors.New("user is not found")
	}

	return u, nil
}

// EmailExists checks if user with specified email exists
func EmailExists(email string) bool {
	id := 0

	err := pg.DB.QueryRow(
		context.Background(),
		"SELECT users.id FROM users WHERE lower(users.email) = $1 LIMIT 1",
		strings.ToLower(email),
	).Scan(&id)

	if err != nil && err.Error() != ignoredSQLError {
		utils.LogError(err.Error())
	}

	if id == 0 {
		return false
	}

	return true
}
