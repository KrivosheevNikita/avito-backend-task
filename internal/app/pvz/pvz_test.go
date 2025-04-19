package pvz

import (
	"errors"
	"pvz-service/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct {
	savePVZErr                     error
	fetchPVZResult                 []models.PVZ
	fetchPVZErr                    error
	fetchReceptionsByPVZResult     map[uuid.UUID][]models.Reception
	fetchReceptionsByPVZErr        error
	fetchProductsByReceptionResult map[uuid.UUID][]models.Product
	fetchProductsByReceptionErr    error
	updateReceptionErr             error
}

func (m *mockStorage) SavePVZ(p *models.PVZ) error {
	p.ID = [16]byte{1}
	p.RegistrationDate = time.Now()
	return m.savePVZErr
}

func (m *mockStorage) FetchPVZ(startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error) {
	return m.fetchPVZResult, m.fetchPVZErr
}

func (m *mockStorage) FetchReceptionsByPVZ(pvzID uuid.UUID) ([]models.Reception, error) {
	return m.fetchReceptionsByPVZResult[pvzID], m.fetchReceptionsByPVZErr
}

func (m *mockStorage) FetchProductsByReception(receptionID uuid.UUID) ([]models.Product, error) {
	return m.fetchProductsByReceptionResult[receptionID], m.fetchProductsByReceptionErr
}

func (m *mockStorage) UpdateReception(r *models.Reception) error {
	return m.updateReceptionErr
}

func TestCreatePVZ_Success(t *testing.T) {
	originalRepo := repo
	defer func() { repo = originalRepo }()

	repo = &mockStorage{}

	pvz := &models.PVZ{City: "Москва"}
	err := CreatePVZ(pvz)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, pvz.ID)
	assert.False(t, pvz.RegistrationDate.IsZero())
}

func TestCreatePVZ_InvalidCity(t *testing.T) {
	pvz := &models.PVZ{City: "123"}
	err := CreatePVZ(pvz)
	assert.Equal(t, models.ErrInvalidCity, err)
}

func TestCreatePVZ_SaveError(t *testing.T) {
	originalRepo := repo
	defer func() { repo = originalRepo }()

	repo = &mockStorage{savePVZErr: errors.New("error")}

	pvz := &models.PVZ{City: "Казань"}
	err := CreatePVZ(pvz)
	assert.EqualError(t, err, "error")
}

func TestGetPVZList_Success(t *testing.T) {
	originalRepo := repo
	defer func() { repo = originalRepo }()

	now := time.Now()
	pvzID := uuid.New()
	receptionID := uuid.New()

	repo = &mockStorage{
		fetchPVZResult: []models.PVZ{
			{
				ID:               pvzID,
				RegistrationDate: now,
				City:             "Москва",
			},
		},
		fetchReceptionsByPVZResult: map[uuid.UUID][]models.Reception{
			pvzID: {
				{
					ID:       receptionID,
					DateTime: now,
					PvzID:    pvzID,
					Status:   "in_progress",
				},
			},
		},
		fetchProductsByReceptionResult: map[uuid.UUID][]models.Product{
			receptionID: {
				{
					ID:          uuid.New(),
					DateTime:    now,
					Type:        "электроника",
					ReceptionID: receptionID,
				},
			},
		},
	}

	start := now.Add(-time.Hour)
	end := now.Add(time.Hour)
	result, err := GetPVZList(&start, &end, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, "Москва", result[0].PVZ.City)
	assert.Equal(t, "электроника", result[0].Receptions[0].Products[0].Type)
}

func TestGetPVZList_EmptyResult(t *testing.T) {
	originalRepo := repo
	defer func() { repo = originalRepo }()

	repo = &mockStorage{
		fetchPVZResult: []models.PVZ{},
	}

	result, err := GetPVZList(nil, nil, 1, 10)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetPVZList_InvalidDateRange(t *testing.T) {
	start := time.Now()
	end := start.Add(-time.Hour)
	_, err := GetPVZList(&start, &end, 1, 10)
	assert.Equal(t, models.ErrInvalidDateRange, err)
}

func TestGetPVZList_FetchPVZError(t *testing.T) {
	originalRepo := repo
	defer func() { repo = originalRepo }()

	repo = &mockStorage{
		fetchPVZErr: errors.New("error"),
	}

	_, err := GetPVZList(nil, nil, 1, 10)
	assert.EqualError(t, err, "error")
}

func TestGetPVZList_FetchReceptionsError(t *testing.T) {
	originalRepo := repo
	defer func() { repo = originalRepo }()

	pvzID := uuid.New()
	repo = &mockStorage{
		fetchPVZResult: []models.PVZ{
			{
				ID:               pvzID,
				RegistrationDate: time.Now(),
				City:             "Москва",
			},
		},
		fetchReceptionsByPVZErr: errors.New("error"),
	}

	_, err := GetPVZList(nil, nil, 1, 10)
	assert.EqualError(t, err, "error")
}

func TestGetPVZList_FetchProductsError(t *testing.T) {
	originalRepo := repo
	defer func() { repo = originalRepo }()

	pvzID := uuid.New()
	receptionID := uuid.New()
	repo = &mockStorage{
		fetchPVZResult: []models.PVZ{
			{
				ID:               pvzID,
				RegistrationDate: time.Now(),
				City:             "Москва",
			},
		},
		fetchReceptionsByPVZResult: map[uuid.UUID][]models.Reception{
			pvzID: {
				{
					ID:       receptionID,
					DateTime: time.Now(),
					PvzID:    pvzID,
					Status:   "in_progress",
				},
			},
		},
		fetchProductsByReceptionErr: errors.New("error"),
	}

	_, err := GetPVZList(nil, nil, 1, 10)
	assert.EqualError(t, err, "error")
}

func TestCloseLastReception_Success(t *testing.T) {
	originalRepo := repo
	defer func() { repo = originalRepo }()

	pvzID := uuid.New()
	receptionID := uuid.New()

	repo = &mockStorage{
		fetchReceptionsByPVZResult: map[uuid.UUID][]models.Reception{
			pvzID: {
				{
					ID:     receptionID,
					PvzID:  pvzID,
					Status: "in_progress",
				},
			},
		},
	}

	result, err := CloseLastReception(pvzID)
	assert.NoError(t, err)
	assert.Equal(t, "close", result.Status)
	assert.Equal(t, receptionID, result.ID)
}

func TestCloseLastReception_NoOpenReception(t *testing.T) {
	originalRepo := repo
	defer func() { repo = originalRepo }()

	pvzID := uuid.New()

	repo = &mockStorage{
		fetchReceptionsByPVZResult: map[uuid.UUID][]models.Reception{
			pvzID: {
				{
					ID:     uuid.New(),
					PvzID:  pvzID,
					Status: "close",
				},
			},
		},
	}

	_, err := CloseLastReception(pvzID)
	assert.Equal(t, models.ErrNoOpenReception, err)
}

func TestCloseLastReception_UpdateError(t *testing.T) {
	originalRepo := repo
	defer func() { repo = originalRepo }()

	pvzID := uuid.New()
	receptionID := uuid.New()

	repo = &mockStorage{
		fetchReceptionsByPVZResult: map[uuid.UUID][]models.Reception{
			pvzID: {
				{
					ID:     receptionID,
					PvzID:  pvzID,
					Status: "in_progress",
				},
			},
		},
		updateReceptionErr: errors.New("error"),
	}

	_, err := CloseLastReception(pvzID)
	assert.EqualError(t, err, "error")
}

func TestCloseLastReception_FetchError(t *testing.T) {
	originalRepo := repo
	defer func() { repo = originalRepo }()

	pvzID := uuid.New()

	repo = &mockStorage{
		fetchReceptionsByPVZErr: errors.New("error"),
	}

	_, err := CloseLastReception(pvzID)
	assert.EqualError(t, err, "error")
}
