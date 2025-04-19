package db

import (
	"fmt"
	"pvz-service/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSavePVZ(t *testing.T) {
	mockDB := new(mockDBExecutor)
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil)

	SetDBExecutor(mockDB)

	newPVZ := &models.PVZ{
		ID:               uuid.New(),
		City:             "Москва",
		RegistrationDate: time.Now(),
	}

	err := SavePVZ(newPVZ)
	assert.Nil(t, err)
	mockDB.AssertExpectations(t)
}

func TestSavePVZ_Error(t *testing.T) {
	mockDB := new(mockDBExecutor)
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, fmt.Errorf("ошибка бд"))

	SetDBExecutor(mockDB)

	newPVZ := &models.PVZ{
		ID:               uuid.New(),
		City:             "Москва",
		RegistrationDate: time.Now(),
	}

	err := SavePVZ(newPVZ)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Ошибка при выполнении запроса на создание ПВЗ: ошибка бд")
	mockDB.AssertExpectations(t)
}

func TestFetchPVZ_Error(t *testing.T) {
	mockDB := new(mockDBExecutor)
	mockDB.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(pgx.Rows(nil), fmt.Errorf("ошибка бд"))

	SetDBExecutor(mockDB)

	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	_, err := FetchPVZ(&startDate, &endDate, 1, 10)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "ошибка бд")
	mockDB.AssertExpectations(t)
}

func TestFetchReceptionsByPVZ_Error(t *testing.T) {
	mockDB := new(mockDBExecutor)
	mockDB.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("ошибка бд"))

	SetDBExecutor(mockDB)

	pvzID := uuid.New()

	_, err := FetchReceptionsByPVZ(pvzID)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "ошибка бд")
	mockDB.AssertExpectations(t)
}

func TestGetPvz_DBError(t *testing.T) {
	mockDB := new(mockDBExecutor)
	SetDBExecutor(mockDB)

	mockDB.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("ошибка бд"))

	result, err := GetPvz()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "ошибка бд")
	mockDB.AssertExpectations(t)
}
