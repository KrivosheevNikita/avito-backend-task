package product

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pvz-service/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddProductHandler_Success(t *testing.T) {
	pvzID := uuid.New()
	receptionID := uuid.New()
	now := time.Now()

	AddProductFunc = func(pvzId uuid.UUID, productType string) (*models.Product, error) {
		return &models.Product{
			ID:          uuid.New(),
			DateTime:    now,
			Type:        productType,
			ReceptionID: receptionID,
		}, nil
	}

	reqBody := AddProductRequest{
		Type:  "электроника",
		PvzId: pvzID.String(),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	AddProductHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp AddProductResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, reqBody.Type, resp.Type)
	assert.Equal(t, receptionID, resp.ReceptionID)
	assert.WithinDuration(t, now, resp.DateTime, time.Second)
}

func TestAddProductHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte(`{invalid json`)))
	rr := httptest.NewRecorder()

	AddProductHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddProductHandler_InvalidUUID(t *testing.T) {
	body := []byte(`{"type": "одежда", "pvzId": "wrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	AddProductHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddProductHandler_NoActiveReception(t *testing.T) {
	AddProductFunc = func(pvzId uuid.UUID, productType string) (*models.Product, error) {
		return nil, models.ErrNoActiveReception
	}

	pvzID := uuid.New()
	body, _ := json.Marshal(AddProductRequest{
		Type:  "обувь",
		PvzId: pvzID.String(),
	})

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	AddProductHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddProductHandler_InvalidType(t *testing.T) {
	AddProductFunc = func(pvzId uuid.UUID, productType string) (*models.Product, error) {
		return nil, models.ErrInvalidTypeProduct
	}

	pvzID := uuid.New()
	body, _ := json.Marshal(AddProductRequest{
		Type:  "123",
		PvzId: pvzID.String(),
	})

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	AddProductHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddProductHandler_InternalError(t *testing.T) {
	AddProductFunc = func(pvzId uuid.UUID, productType string) (*models.Product, error) {
		return nil, errors.New("Ошибка бд")
	}

	pvzID := uuid.New()
	body, _ := json.Marshal(AddProductRequest{
		Type:  "одежда",
		PvzId: pvzID.String(),
	})

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	AddProductHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
