package reception_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pvz-service/internal/models"
	handler "pvz-service/internal/transport/handlers/reception"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func overrideCreateReception(fn func(pvzID uuid.UUID) (*models.Reception, error)) {
	handler.CreateReceptionFunc = fn
}

func TestCreateReceptionHandler_Success(t *testing.T) {
	testUUID := uuid.New()
	expected := &models.Reception{
		ID:       uuid.New(),
		PvzID:    testUUID,
		DateTime: time.Now().UTC(),
		Status:   "in_progress",
	}

	overrideCreateReception(func(pvzID uuid.UUID) (*models.Reception, error) {
		assert.Equal(t, testUUID, pvzID)
		return expected, nil
	})

	body, _ := json.Marshal(handler.CreateReceptionRequest{PvzID: testUUID.String()})
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateReceptionHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response handler.CreateReceptionResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, response.ID)
	assert.Equal(t, expected.PvzID, response.PvzID)
	assert.Equal(t, expected.Status, "in_progress")
	assert.WithinDuration(t, expected.DateTime, response.DateTime, time.Second)
}

func TestCreateReceptionHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader([]byte("invalid-json")))
	w := httptest.NewRecorder()

	handler.CreateReceptionHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestCreateReceptionHandler_InvalidUUID(t *testing.T) {
	body := []byte(`{"pvz_id": "wrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateReceptionHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestCreateReceptionHandler_BusinessError(t *testing.T) {
	testUUID := uuid.New()

	overrideCreateReception(func(pvzID uuid.UUID) (*models.Reception, error) {
		return nil, models.ErrReceptionInProgress
	})

	body, _ := json.Marshal(handler.CreateReceptionRequest{PvzID: testUUID.String()})
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateReceptionHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}
