package pvz

import (
	"context"
	"testing"
	"time"

	"pvz-service/api/grpc/pvz/pvz_v1"
	"pvz-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetPvz() ([]models.PVZ, error) {
	args := m.Called()
	return args.Get(0).([]models.PVZ), args.Error(1)
}

func TestGetPVZList_Success(t *testing.T) {
	mockDB := new(MockDB)
	GetPvzFunc = mockDB.GetPvz

	pvzs := []models.PVZ{
		{
			ID:               uuid.New(),
			RegistrationDate: time.Now().UTC(),
			City:             "Москва",
		},
		{
			ID:               uuid.New(),
			RegistrationDate: time.Now().UTC(),
			City:             "Санкт-Петербург",
		},
	}

	mockDB.On("GetPvz").Return(pvzs, nil)

	server := NewPVZServer()
	req := &pvz_v1.GetPVZListRequest{}
	ctx := context.Background()

	resp, err := server.GetPVZList(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, resp.Pvzs, 2)
	assert.Equal(t, "Москва", resp.Pvzs[0].City)
	assert.Equal(t, "Санкт-Петербург", resp.Pvzs[1].City)
	mockDB.AssertExpectations(t)
}
func TestGetPVZList_EmptyList(t *testing.T) {
	mockDB := new(MockDB)
	GetPvzFunc = mockDB.GetPvz

	mockDB.On("GetPvz").Return([]models.PVZ{}, nil)

	server := NewPVZServer()
	req := &pvz_v1.GetPVZListRequest{}
	ctx := context.Background()

	resp, err := server.GetPVZList(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, resp.Pvzs, 0)
	mockDB.AssertExpectations(t)
}
