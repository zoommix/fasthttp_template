package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	jwt "github.com/zoommix/fasthttp_template/utils"
)

// TestGenerateToken ...
func TestGenerateToken(t *testing.T) {
	jwtWrapper := jwt.JwtWrapper{
		SecretKey:       "verysecretkey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	generatedToken, err := jwtWrapper.GenerateToken(121)
	assert.NoError(t, err)

	os.Setenv("testToken", generatedToken)
}

// TestValidateToken ...
func TestValidateToken(t *testing.T) {
	encodedToken := os.Getenv("testToken")

	jwtWrapper := jwt.JwtWrapper{
		SecretKey: "verysecretkey",
		Issuer:    "AuthService",
	}

	claims, err := jwtWrapper.ValidateToken(encodedToken)
	assert.NoError(t, err)

	assert.Equal(t, 121, claims.UserID)
	assert.Equal(t, "AuthService", claims.Issuer)
}
