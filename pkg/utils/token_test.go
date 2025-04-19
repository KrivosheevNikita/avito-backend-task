package utils

import (
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "key")
	os.Setenv("JWT_EXPIRY_HOURS", "1")

	token, err := GenerateToken("moderator")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "moderator", claims["role"])
}

func TestGenerateToken_InvalidExpiry(t *testing.T) {
	os.Setenv("JWT_SECRET", "key")
	os.Setenv("JWT_EXPIRY_HOURS", "invalid")

	token, err := GenerateToken("moderator")
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestParseToken_Success(t *testing.T) {

	os.Setenv("JWT_SECRET", "key")
	os.Setenv("JWT_EXPIRY_HOURS", "1")

	token, err := GenerateToken("moderator")
	assert.NoError(t, err)

	claims, err := ParseToken(token)
	assert.NoError(t, err)

	assert.Equal(t, "moderator", claims["role"])
}

func TestParseToken_InvalidToken(t *testing.T) {

	os.Setenv("JWT_SECRET", "key")

	invalidToken := "invalid"

	claims, err := ParseToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestParseToken_UnexpectedSigningMethod(t *testing.T) {

	os.Setenv("JWT_SECRET", "key")

	claims := jwt.MapClaims{
		"role": "moderator",
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString([]byte("key"))

	parsedClaims, err := ParseToken(signedToken)

	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
}
