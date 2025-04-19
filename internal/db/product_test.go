package db

import (
	"fmt"
	"pvz-service/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddProductToReception(t *testing.T) {
	mockDB := new(mockDBExecutor)
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil)

	SetDBExecutor(mockDB)

	product := &models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        "электроника",
		ReceptionID: uuid.New(),
	}
	err := AddProductToReception(product)
	assert.Nil(t, err)
	mockDB.AssertExpectations(t)
}

func TestAddProduct_Error(t *testing.T) {
	mockDB := new(mockDBExecutor)
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, fmt.Errorf("ошибка бд"))

	SetDBExecutor(mockDB)

	product := &models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        "электроника",
		ReceptionID: uuid.New(),
	}

	err := AddProductToReception(product)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Ошибка при выполнении запроса на добавление товара: ошибка бд")
	mockDB.AssertExpectations(t)
}

func TestDeleteLastProductByReceptionID_Error(t *testing.T) {
	mockDB := new(mockDBExecutor)
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, fmt.Errorf("ошибка бд"))

	SetDBExecutor(mockDB)

	receptionID := uuid.New()

	err := DeleteLastProductByReceptionID(receptionID)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Ошибка при удалении последнего товара: ошибка бд")
	mockDB.AssertExpectations(t)
}

func TestDeleteLastProductByReceptionID_NoRowsAffected(t *testing.T) {
	mockDB := new(mockDBExecutor)
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil)

	SetDBExecutor(mockDB)

	receptionID := uuid.New()

	err := DeleteLastProductByReceptionID(receptionID)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "Нет товара для удаления")

	mockDB.AssertExpectations(t)
}
