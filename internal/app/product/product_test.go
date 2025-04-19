package product

import (
	"errors"
	"pvz-service/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct {
	getReceptionResult *models.Reception
	getReceptionErr    error
	addProductErr      error
	addedProduct       *models.Product

	deleteErr        error
	deleteCalledWith uuid.UUID
}

func (m *mockStorage) GetActiveReceptionByPVZ(pvzID uuid.UUID) (*models.Reception, error) {
	return m.getReceptionResult, m.getReceptionErr
}

func (m *mockStorage) AddProductToReception(p *models.Product) error {
	m.addedProduct = p
	return m.addProductErr
}

func (m *mockStorage) DeleteLastProductByReceptionID(receptionID uuid.UUID) error {
	m.deleteCalledWith = receptionID
	return m.deleteErr
}

func TestAddProduct_Success(t *testing.T) {
	pvzID := uuid.New()
	receptionID := uuid.New()

	mock := &mockStorage{
		getReceptionResult: &models.Reception{ID: receptionID},
	}
	repo = mock

	productType := "электроника"
	product, err := AddProduct(pvzID, productType)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, productType, product.Type)
	assert.Equal(t, receptionID, product.ReceptionID)
	assert.WithinDuration(t, time.Now().UTC(), product.DateTime, time.Second*1)
	assert.Equal(t, mock.addedProduct.ID, product.ID)
}

func TestAddProduct_NoActiveReception(t *testing.T) {
	pvzID := uuid.New()

	mock := &mockStorage{
		getReceptionErr: models.ErrNoActiveReception,
	}
	repo = mock

	product, err := AddProduct(pvzID, "обувь")

	assert.Nil(t, product)
	assert.ErrorIs(t, err, models.ErrNoActiveReception)
}

func TestAddProduct_GetReceptionError(t *testing.T) {
	pvzID := uuid.New()

	mock := &mockStorage{
		getReceptionErr: errors.New("error"),
	}
	repo = mock

	product, err := AddProduct(pvzID, "одежда")

	assert.Nil(t, product)
	assert.EqualError(t, err, "error")
}

func TestAddProduct_SaveError(t *testing.T) {
	pvzID := uuid.New()
	receptionID := uuid.New()

	mock := &mockStorage{
		getReceptionResult: &models.Reception{ID: receptionID},
		addProductErr:      errors.New("error"),
	}
	repo = mock

	product, err := AddProduct(pvzID, "электроника")

	assert.Nil(t, product)
	assert.EqualError(t, err, "error")
}

func TestAddProduct_InvalidType(t *testing.T) {
	pvzID := uuid.New()
	mock := &mockStorage{}
	repo = mock

	product, err := AddProduct(pvzID, "123")

	assert.Nil(t, product)
	assert.ErrorIs(t, err, models.ErrInvalidTypeProduct)
}

func TestDeleteLastProduct_Success(t *testing.T) {
	pvzID := uuid.New()
	receptionID := uuid.New()

	mock := &mockStorage{
		getReceptionResult: &models.Reception{ID: receptionID},
	}
	repo = mock

	err := DeleteLastProduct(pvzID)

	assert.NoError(t, err)
	assert.Equal(t, receptionID, mock.deleteCalledWith)
}

func TestDeleteLastProduct_NoActiveReception(t *testing.T) {
	pvzID := uuid.New()

	mock := &mockStorage{
		getReceptionErr: models.ErrNoActiveReception,
	}
	repo = mock

	err := DeleteLastProduct(pvzID)

	assert.ErrorIs(t, err, models.ErrNoActiveReception)
}

func TestDeleteLastProduct_GetReceptionError(t *testing.T) {
	pvzID := uuid.New()

	mock := &mockStorage{
		getReceptionErr: errors.New("error"),
	}
	repo = mock

	err := DeleteLastProduct(pvzID)

	assert.EqualError(t, err, "error")
}

func TestDeleteLastProduct_DeleteError(t *testing.T) {
	pvzID := uuid.New()
	receptionID := uuid.New()

	mock := &mockStorage{
		getReceptionResult: &models.Reception{ID: receptionID},
		deleteErr:          errors.New("error"),
	}
	repo = mock

	err := DeleteLastProduct(pvzID)

	assert.EqualError(t, err, "error")
}
