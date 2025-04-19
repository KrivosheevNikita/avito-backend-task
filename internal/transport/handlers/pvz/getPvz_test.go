package pvz

import (
	"net/http"
	"net/http/httptest"
	"pvz-service/internal/app/pvz"
	"pvz-service/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct{}

func (m *mockRepo) GetPVZList(startDate, endDate *time.Time, page, limit int) ([]pvz.PVZWithReceptions, error) {
	return []pvz.PVZWithReceptions{
		{
			PVZ: models.PVZ{
				ID:               uuid.New(),
				City:             "Москва",
				RegistrationDate: time.Now(),
			},
			Receptions: []pvz.ReceptionGroup{
				{
					Reception: models.Reception{
						ID: uuid.New(),
					},
					Products: []models.Product{
						{
							ID:   uuid.New(),
							Type: "электроника",
						},
						{
							ID:   uuid.New(),
							Type: "одежда",
						},
					},
				},
			},
		},
	}, nil
}
func TestGetPVZHandler_Success(t *testing.T) {
	GetPVZListFunc = (&mockRepo{}).GetPVZList

	req, err := http.NewRequest("GET", "/pvz?startDate=2024-01-01T00:00:00Z&endDate=2025-01-01T00:00:00Z&page=1&limit=10", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPVZHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetPVZHandler_InvalidDateFormat(t *testing.T) {
	GetPVZListFunc = (&mockRepo{}).GetPVZList

	req, err := http.NewRequest("GET", "/pvz?startDate=invalid-date&endDate=invalid-date", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPVZHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetPVZHandler_EmptyResult(t *testing.T) {
	GetPVZListFunc = func(startDate, endDate *time.Time, page, limit int) ([]pvz.PVZWithReceptions, error) {
		return []pvz.PVZWithReceptions{}, nil
	}

	req, err := http.NewRequest("GET", "/pvz?startDate=2025-01-01T00:00:00Z&endDate=2025-01-02T00:00:00Z&page=1&limit=10", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPVZHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetPVZHandler_InvalidPage(t *testing.T) {
	GetPVZListFunc = (&mockRepo{}).GetPVZList

	tests := []string{
		"/pvz?page=0",
		"/pvz?page=abc",
	}

	for _, url := range tests {
		req, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetPVZHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	}
}

func TestGetPVZHandler_InvalidLimit(t *testing.T) {
	GetPVZListFunc = (&mockRepo{}).GetPVZList

	tests := []string{
		"/pvz?limit=0",
		"/pvz?limit=1000",
		"/pvz?limit=abc",
	}

	for _, url := range tests {
		req, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetPVZHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	}
}
