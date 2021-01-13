package users

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	pg "github.com/zoommix/fasthttp_template/store"
	log "github.com/zoommix/fasthttp_template/utils"
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
		log.LogInfo(fmt.Sprintf("Unable to create user with email=%+v", u))

		return nil, errors.New(strings.Replace(err.Error(), `"`, `'`, -1))
	}

	log.LogInfo(fmt.Sprintf("Created user with ID=%d, email=%s", u.ID, u.Email))

	return u, nil
}

// EncodePassword ...
func EncodePassword(p string) string {
	sha := sha256.Sum256([]byte(p))
	return base64.StdEncoding.EncodeToString(sha[:])
}
