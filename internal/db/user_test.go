package db

import (
	"fmt"
	"pvz-service/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Success(t *testing.T) {
	mockDB := new(mockDBExecutor)
	SetDBExecutor(mockDB)

	user := &models.User{
		ID:    uuid.New(),
		Email: "email",
		Role:  "moderator",
	}

	mockDB.On("QueryRow", mock.Anything, mock.Anything, user.Email).
		Return(mockDBRow{scanFn: func(dest ...interface{}) error {
			*(dest[0].(*bool)) = false
			return nil
		}})

	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.NewCommandTag("INSERT 1"), nil)

	err := CreateUser(user, "password")
	assert.Nil(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
	mockDB.AssertExpectations(t)
}

func TestCreateUser_ErrorOnInsert(t *testing.T) {
	mockDB := new(mockDBExecutor)
	SetDBExecutor(mockDB)

	user := &models.User{
		ID:    uuid.New(),
		Email: "email",
		Role:  "moderator",
	}

	mockDB.On("QueryRow", mock.Anything, mock.Anything, user.Email).
		Return(mockDBRow{scanFn: func(dest ...interface{}) error {
			*(dest[0].(*bool)) = false
			return nil
		}})

	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, fmt.Errorf("Ошибка при сохранении пользователя в базу данных:"))

	err := CreateUser(user, "password")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Ошибка при сохранении пользователя в базу данных:")
	mockDB.AssertExpectations(t)
}

func TestGetUserByEmail_Success(t *testing.T) {
	mockDB := new(mockDBExecutor)
	SetDBExecutor(mockDB)

	testID := uuid.New()
	email := "email"
	role := "moderator"
	passwordHash := "hash"

	mockDB.On("QueryRow", mock.Anything, mock.Anything, email).
		Return(mockDBRow{scanFn: func(dest ...interface{}) error {
			*(dest[0].(*uuid.UUID)) = testID
			*(dest[1].(*string)) = email
			*(dest[2].(*string)) = role
			*(dest[3].(*string)) = passwordHash
			return nil
		}})

	user, pass, err := GetUserByEmail(email)
	assert.Nil(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, role, user.Role)
	assert.Equal(t, passwordHash, pass)
	mockDB.AssertExpectations(t)
}

func TestGetUserByEmail_Error(t *testing.T) {
	mockDB := new(mockDBExecutor)
	SetDBExecutor(mockDB)

	email := "email"

	mockDB.On("QueryRow", mock.Anything, mock.Anything, email).
		Return(mockDBRow{scanFn: func(dest ...interface{}) error {
			return fmt.Errorf("Пользователь не найден")
		}})

	user, pass, err := GetUserByEmail(email)

	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Empty(t, pass)
	mockDB.AssertExpectations(t)
}
