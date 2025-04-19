package product

import (
	"pvz-service/internal/db"
	"pvz-service/internal/metrics"
	"pvz-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type storage interface {
	GetActiveReceptionByPVZ(pvzID uuid.UUID) (*models.Reception, error)
	AddProductToReception(product *models.Product) error
	DeleteLastProductByReceptionID(receptionID uuid.UUID) error
}

type dbAdapter struct{}

func (dbAdapter) GetActiveReceptionByPVZ(pvzID uuid.UUID) (*models.Reception, error) {
	return db.GetActiveReceptionByPVZ(pvzID)
}

func (dbAdapter) AddProductToReception(product *models.Product) error {
	return db.AddProductToReception(product)
}

func (dbAdapter) DeleteLastProductByReceptionID(receptionID uuid.UUID) error {
	return db.DeleteLastProductByReceptionID(receptionID)
}

var repo storage = dbAdapter{}

func AddProduct(pvzId uuid.UUID, productType string) (*models.Product, error) {
	if productType != "электроника" && productType != "обувь" && productType != "одежда" {
		return nil, models.ErrInvalidTypeProduct
	}
	reception, err := repo.GetActiveReceptionByPVZ(pvzId)
	if err != nil {
		return nil, err
	}
	if reception == nil {
		return nil, models.ErrNoActiveReception
	}

	product := &models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        productType,
		ReceptionID: reception.ID,
	}

	err = repo.AddProductToReception(product)
	if err != nil {
		return nil, err
	}

	metrics.ProductsAdded.Inc()

	return product, nil
}

func DeleteLastProduct(pvzID uuid.UUID) error {

	reception, err := repo.GetActiveReceptionByPVZ(pvzID)
	if err != nil {
		return err
	}

	if reception == nil {
		return models.ErrNoActiveReception
	}

	err = repo.DeleteLastProductByReceptionID(reception.ID)
	if err != nil {
		return err
	}

	return nil
}
