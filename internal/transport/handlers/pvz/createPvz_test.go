package pvz_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pvz-service/internal/models"
	handler "pvz-service/internal/transport/handlers/pvz"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func overrideCreatePVZ(fn func(p *models.PVZ) error) {
	handler.CreatePVZFunc = fn
}

func TestCreatePVZHandler_Success(t *testing.T) {
	overrideCreatePVZ(func(p *models.PVZ) error {
		p.ID = uuid.New()
		p.RegistrationDate = time.Now().UTC()
		return nil
	})

	body, _ := json.Marshal(handler.CreatePVZRequest{City: "Москва"})
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreatePVZHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response handler.CreatePVZResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Москва", response.City)
	assert.NotEmpty(t, response.ID)
	assert.WithinDuration(t, time.Now().UTC(), response.RegistrationDate, time.Second)
}

func TestCreatePVZHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader([]byte("invalid-json")))
	w := httptest.NewRecorder()

	handler.CreatePVZHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestCreatePVZHandler_InvalidCity(t *testing.T) {
	overrideCreatePVZ(func(p *models.PVZ) error {
		return models.ErrInvalidCity
	})

	body, _ := json.Marshal(handler.CreatePVZRequest{City: "123"})
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreatePVZHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}
