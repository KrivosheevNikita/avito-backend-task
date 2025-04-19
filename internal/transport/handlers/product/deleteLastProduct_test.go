package product_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"pvz-service/internal/models"
	handler "pvz-service/internal/transport/handlers/product"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func overrideDeleteLastProduct(fn func(uuid.UUID) error) {
	handler.DeleteLastProductFunc = fn
}

func setURLVar(r *http.Request, key, val string) *http.Request {
	return mux.SetURLVars(r, map[string]string{key: val})
}

func TestDeleteLastProductHandler_Success(t *testing.T) {
	pvzID := uuid.New()

	overrideDeleteLastProduct(func(id uuid.UUID) error {
		assert.Equal(t, pvzID, id)
		return nil
	})

	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product", nil)
	req = setURLVar(req, "pvzId", pvzID.String())

	w := httptest.NewRecorder()
	handler.DeleteLastProductHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestDeleteLastProductHandler_InvalidUUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/pvz/invalid-uuid/delete_last_product", nil)
	req = setURLVar(req, "pvzId", "invalid-uuid")

	w := httptest.NewRecorder()
	handler.DeleteLastProductHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestDeleteLastProductHandler_NoActiveReception(t *testing.T) {
	overrideDeleteLastProduct(func(id uuid.UUID) error {
		return models.ErrNoActiveReception
	})

	id := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/pvz/"+id.String()+"/delete_last_product", nil)
	req = setURLVar(req, "pvzId", id.String())

	w := httptest.NewRecorder()
	handler.DeleteLastProductHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestDeleteLastProductHandler_ServerError(t *testing.T) {
	overrideDeleteLastProduct(func(id uuid.UUID) error {
		return errors.New("Ошибка бд")
	})

	id := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/pvz/"+id.String()+"/delete_last_product", nil)
	req = setURLVar(req, "pvzId", id.String())

	w := httptest.NewRecorder()
	handler.DeleteLastProductHandler(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestDeleteLastProductHandler_MissingPvzID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/pvz/id/delete_last_product", nil)

	w := httptest.NewRecorder()
	handler.DeleteLastProductHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}
