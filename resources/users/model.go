package users

import (
	"os"

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
