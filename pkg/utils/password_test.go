package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "123"

	hashedPassword, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEqual(t, password, hashedPassword)
	assert.NotEmpty(t, hashedPassword)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "123"

	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)

	isValid := CheckPasswordHash(password, hashedPassword)
	assert.True(t, isValid, "")

	incorrectPassword := "wrong"
	isValid = CheckPasswordHash(incorrectPassword, hashedPassword)
	assert.False(t, isValid, "")
}
