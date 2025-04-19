package reception

import (
	"pvz-service/internal/db"
	"pvz-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type storage interface {
	HasActiveReception(pvzID uuid.UUID) (bool, error)
	SaveReception(r *models.Reception) error
}

type dbAdapter struct{}

func (dbAdapter) HasActiveReception(pvzID uuid.UUID) (bool, error) {
	return db.HasActiveReception(pvzID)
}

func (dbAdapter) SaveReception(r *models.Reception) error {
	return db.SaveReception(r)
}

var repo storage = dbAdapter{}

func CreateReception(pvzID uuid.UUID) (*models.Reception, error) {
	inProgress, err := repo.HasActiveReception(pvzID)
	if err != nil {
		return nil, err
	}
	if inProgress {
		return nil, models.ErrReceptionInProgress
	}

	newReception := &models.Reception{
		ID:       uuid.New(),
		PvzID:    pvzID,
		DateTime: time.Now(),
		Status:   "in_progress",
	}

	if err := repo.SaveReception(newReception); err != nil {
		return nil, err
	}

	return newReception, nil
}
