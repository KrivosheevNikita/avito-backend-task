package pvz

import (
	"pvz-service/internal/db"
	"pvz-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type storage interface {
	SavePVZ(p *models.PVZ) error
	FetchPVZ(startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error)
	FetchReceptionsByPVZ(pvzID uuid.UUID) ([]models.Reception, error)
	FetchProductsByReception(receptionID uuid.UUID) ([]models.Product, error)
	UpdateReception(reception *models.Reception) error
}

type dbAdapter struct{}

func (dbAdapter) SavePVZ(p *models.PVZ) error {
	return db.SavePVZ(p)
}

func (dbAdapter) FetchPVZ(startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error) {
	return db.FetchPVZ(startDate, endDate, page, limit)
}

func (dbAdapter) FetchReceptionsByPVZ(pvzID uuid.UUID) ([]models.Reception, error) {
	return db.FetchReceptionsByPVZ(pvzID)
}

func (dbAdapter) FetchProductsByReception(receptionID uuid.UUID) ([]models.Product, error) {
	return db.FetchProductsByReception(receptionID)
}

func (dbAdapter) UpdateReception(r *models.Reception) error {
	return db.UpdateReception(r)
}

var repo storage = dbAdapter{}

type PVZWithReceptions struct {
	PVZ        models.PVZ       `json:"pvz"`
	Receptions []ReceptionGroup `json:"receptions"`
}

type ReceptionGroup struct {
	Reception models.Reception `json:"reception"`
	Products  []models.Product `json:"products"`
}

func CreatePVZ(newPVZ *models.PVZ) error {
	if newPVZ.City != "Москва" && newPVZ.City != "Санкт-Петербург" && newPVZ.City != "Казань" {
		return models.ErrInvalidCity
	}
	newPVZ.ID = uuid.New()
	newPVZ.RegistrationDate = time.Now()
	return repo.SavePVZ(newPVZ)
}

func GetPVZList(startDate, endDate *time.Time, page, limit int) ([]PVZWithReceptions, error) {
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return nil, models.ErrInvalidDateRange
	}

	pvzs, err := repo.FetchPVZ(startDate, endDate, page, limit)
	if err != nil {
		return nil, err
	}

	var result []PVZWithReceptions
	for _, pvz := range pvzs {
		receptions, err := repo.FetchReceptionsByPVZ(pvz.ID)
		if err != nil {
			return nil, err
		}

		var grouped []ReceptionGroup
		for _, r := range receptions {
			products, err := repo.FetchProductsByReception(r.ID)
			if err != nil {
				return nil, err
			}
			grouped = append(grouped, ReceptionGroup{
				Reception: r,
				Products:  products,
			})
		}

		result = append(result, PVZWithReceptions{
			PVZ:        pvz,
			Receptions: grouped,
		})
	}

	return result, nil
}

func CloseLastReception(pvzId uuid.UUID) (*models.Reception, error) {

	receptions, err := repo.FetchReceptionsByPVZ(pvzId)
	if err != nil {
		return nil, err
	}

	var lastOpenReception *models.Reception
	for _, reception := range receptions {
		if reception.Status == "in_progress" {
			lastOpenReception = &reception
			break
		}
	}

	if lastOpenReception == nil {
		return nil, models.ErrNoOpenReception
	}

	if lastOpenReception.Status == "close" {
		return nil, models.ErrReceptionAlreadyClosed
	}

	lastOpenReception.Status = "close"

	err = repo.UpdateReception(lastOpenReception)
	if err != nil {
		return nil, err
	}

	return lastOpenReception, nil
}
