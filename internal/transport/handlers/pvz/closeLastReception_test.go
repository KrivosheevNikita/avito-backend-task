package pvz_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pvz-service/internal/models"
	handler "pvz-service/internal/transport/handlers/pvz"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func overrideCloseLastReception(fn func(uuid.UUID) (*models.Reception, error)) {
	handler.CloseLastReceptionFunc = fn
}

func TestCloseLastReceptionHandler_Success(t *testing.T) {
	expectedID := uuid.New()
	expectedTime := time.Now().UTC()

	overrideCloseLastReception(func(id uuid.UUID) (*models.Reception, error) {
		return &models.Reception{
			ID:       expectedID,
			DateTime: expectedTime,
			PvzID:    id,
			Status:   "close",
		}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/pvz/"+expectedID.String()+"/close_last_reception", nil)
	req = setURLVar(req, "pvzId", expectedID.String())

	w := httptest.NewRecorder()
	handler.CloseLastReceptionHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response models.Reception
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "close", response.Status)
	assert.Equal(t, expectedID, response.ID)
	assert.WithinDuration(t, expectedTime, response.DateTime, time.Second)
}

func TestCloseLastReceptionHandler_InvalidUUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/pvz/wrong/close_last_reception", nil)
	req = setURLVar(req, "pvzId", "wrong")

	w := httptest.NewRecorder()
	handler.CloseLastReceptionHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestCloseLastReceptionHandler_NoOpenReception(t *testing.T) {
	overrideCloseLastReception(func(id uuid.UUID) (*models.Reception, error) {
		return nil, models.ErrNoOpenReception
	})

	id := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/pvz/"+id.String()+"/close_last_reception", nil)
	req = setURLVar(req, "pvzId", id.String())

	w := httptest.NewRecorder()
	handler.CloseLastReceptionHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestCloseLastReceptionHandler_AlreadyClosed(t *testing.T) {
	overrideCloseLastReception(func(id uuid.UUID) (*models.Reception, error) {
		return nil, models.ErrReceptionAlreadyClosed
	})

	id := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/pvz/"+id.String()+"/close_last_reception", nil)
	req = setURLVar(req, "pvzId", id.String())

	w := httptest.NewRecorder()
	handler.CloseLastReceptionHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestCloseLastReceptionHandler_ServerError(t *testing.T) {
	overrideCloseLastReception(func(id uuid.UUID) (*models.Reception, error) {
		return nil, errors.New("Ошибка бд")
	})

	id := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/pvz/"+id.String()+"/close_last_reception", nil)
	req = setURLVar(req, "pvzId", id.String())

	w := httptest.NewRecorder()
	handler.CloseLastReceptionHandler(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func setURLVar(r *http.Request, key, val string) *http.Request {
	return mux.SetURLVars(r, map[string]string{key: val})
}
