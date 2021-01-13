package users

import (
	"context"
	"errors"
	"os"
	"strings"

	pg "github.com/zoommix/fasthttp_template/store"
	"github.com/zoommix/fasthttp_template/utils"
)

// User ...
type User struct {
	ID             int
	FirstName      string
	LastName       string
	Email          string
	Password       string
	PasswordDigest string
	Token          string
	CreatedAt      int
}

const (
	jwtSaltKey      = "JWT_SECRET"
	ignoredSQLError = "no rows in result set"
)

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

// Reload returns pointer to recent user data
func (u *User) Reload() *User {
	u, err := FindUser(u.ID)

	if err != nil {
		utils.LogError("Unable to reload user")
	}

	return u
}

// GenerateJWT returns JWT token
func (u *User) GenerateJWT() (err error) {
	jwtSaltENV := os.Getenv(jwtSaltKey)
	jwtSalt := "verysecretkey"

	if len(jwtSaltENV) != 0 {
		jwtSalt = jwtSaltENV
	}

	jwtWrapper := utils.JwtWrapper{
		SecretKey:       jwtSalt,
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	generatedToken, err := jwtWrapper.GenerateToken(u.ID)

	if err != nil {
		return
	}

	u.Token = generatedToken

	return
}

// DecodeJWT ...
func DecodeJWT(encodedToken string) (u *User, err error) {
	u = &User{}

	jwtWrapper := utils.JwtWrapper{
		SecretKey: "verysecretkey",
		Issuer:    "AuthService",
	}

	claims, err := jwtWrapper.ValidateToken(encodedToken)

	if err != nil {
		return nil, err
	}

	u.ID = claims.UserID

	return u, nil
}
