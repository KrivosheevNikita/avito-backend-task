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

func TestSaveReception_Success(t *testing.T) {
	mockDB := new(mockDBExecutor)
	SetDBExecutor(mockDB)

	reception := &models.Reception{
		ID:       uuid.New(),
		DateTime: time.Now(),
		PvzID:    uuid.New(),
		Status:   "in_progress",
	}

	mockDB.On("QueryRow", mock.Anything, mock.Anything, reception.ID, reception.DateTime, reception.PvzID, reception.Status).
		Return(mockDBRow{
			scanFn: func(dest ...interface{}) error {
				dest[0] = reception.ID
				dest[1] = reception.DateTime
				dest[2] = reception.PvzID
				dest[3] = reception.Status
				return nil
			},
		})

	err := SaveReception(reception)
	assert.Nil(t, err)
	mockDB.AssertExpectations(t)
}

func TestGetActiveReceptionByPVZ_Success(t *testing.T) {
	mockDB := new(mockDBExecutor)
	SetDBExecutor(mockDB)

	receptionID := uuid.New()
	pvzID := uuid.New()
	date := time.Now()
	status := "in_progress"

	mockDB.On("QueryRow", mock.Anything, mock.Anything, pvzID).
		Return(mockDBRow{
			scanFn: func(dest ...interface{}) error {

				*dest[0].(*uuid.UUID) = receptionID
				*dest[1].(*time.Time) = date
				*dest[2].(*uuid.UUID) = pvzID
				*dest[3].(*string) = status

				return nil
			},
		})

	reception, err := GetActiveReceptionByPVZ(pvzID)

	assert.Nil(t, err)
	assert.NotNil(t, reception)
	assert.Equal(t, receptionID, reception.ID)
	assert.Equal(t, pvzID, reception.PvzID)
	assert.Equal(t, status, reception.Status)
	mockDB.AssertExpectations(t)
}

func TestHasActiveReception_Exists(t *testing.T) {
	mockDB := new(mockDBExecutor)
	SetDBExecutor(mockDB)

	pvzID := uuid.New()

	mockDB.On("QueryRow", mock.Anything, mock.Anything, pvzID).
		Once().
		Return(mockDBRow{
			scanFn: func(dest ...interface{}) error {
				if len(dest) != 1 {
					return fmt.Errorf("Ожидался 1 аргумент вместо %d", len(dest))
				}
				*dest[0].(*bool) = true
				return nil
			},
		})
	mockDB.On("QueryRow", mock.Anything, mock.Anything, pvzID).
		Once().
		Return(mockDBRow{
			scanFn: func(dest ...interface{}) error {
				if len(dest) != 1 {
					return fmt.Errorf("Ожидался 1 аргумент вместо %d", len(dest))
				}
				*dest[0].(*bool) = true
				return nil
			},
		})

	has, err := HasActiveReception(pvzID)

	assert.Nil(t, err)
	assert.True(t, has)
	mockDB.AssertExpectations(t)
}

func TestUpdateReception(t *testing.T) {
	mockDB := new(mockDBExecutor)
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil)

	SetDBExecutor(mockDB)

	reception := &models.Reception{
		ID:     uuid.New(),
		Status: "close",
	}

	err := UpdateReception(reception)
	assert.Nil(t, err)
	mockDB.AssertExpectations(t)
}
