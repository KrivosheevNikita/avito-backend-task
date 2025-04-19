package reception

import (
	"errors"
	"pvz-service/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	hasActive      bool
	hasActiveErr   error
	saveErr        error
	savedReception *models.Reception
}

func (m *mockRepo) HasActiveReception(pvzID uuid.UUID) (bool, error) {
	return m.hasActive, m.hasActiveErr
}

func (m *mockRepo) SaveReception(r *models.Reception) error {
	m.savedReception = r
	return m.saveErr
}

func TestCreateReception_Success(t *testing.T) {
	mock := &mockRepo{hasActive: false}
	repo = mock
	pvzID := uuid.New()

	r, err := CreateReception(pvzID)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, pvzID, r.PvzID)
	assert.Equal(t, "in_progress", r.Status)
	assert.WithinDuration(t, time.Now().UTC(), r.DateTime, time.Second*1)
}

func TestCreateReception_AlreadyInProgress(t *testing.T) {
	mock := &mockRepo{hasActive: true}
	repo = mock
	pvzID := uuid.New()

	r, err := CreateReception(pvzID)
	assert.Nil(t, r)
	assert.Equal(t, models.ErrReceptionInProgress, err)
}

func TestCreateReception_HasActiveError(t *testing.T) {
	dbErr := errors.New("db error")
	mock := &mockRepo{hasActiveErr: dbErr}
	repo = mock

	r, err := CreateReception(uuid.New())
	assert.Nil(t, r)
	assert.ErrorIs(t, err, dbErr)
}

func TestCreateReception_SaveError(t *testing.T) {
	saveErr := errors.New("save failed")
	mock := &mockRepo{hasActive: false, saveErr: saveErr}
	repo = mock

	r, err := CreateReception(uuid.New())
	assert.Nil(t, r)
	assert.ErrorIs(t, err, saveErr)
}
