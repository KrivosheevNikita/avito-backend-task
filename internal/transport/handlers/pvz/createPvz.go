package pvz

import (
	"encoding/json"
	"net/http"
	"pvz-service/internal/app/pvz"
	"pvz-service/internal/models"
	"pvz-service/pkg/handlerutils"
	"pvz-service/pkg/logger"

	"time"

	"github.com/google/uuid"
)

type CreatePVZRequest struct {
	City string `json:"city"`
}

type CreatePVZResponse struct {
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	City             string    `json:"city"`
}

var CreatePVZFunc = pvz.CreatePVZ

func CreatePVZHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogRequest(r, "Получен запрос на создание ПВЗ")

	var req CreatePVZRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlerutils.WriteError(w, models.ErrBadRequest)
		return
	}

	logger.Info("Создание нового ПВЗ", map[string]interface{}{
		"city": req.City,
	})

	newPVZ := &models.PVZ{
		City: req.City,
	}

	if err := CreatePVZFunc(newPVZ); err != nil {
		handlerutils.WriteError(w, err)
		return
	}

	logger.LogPVZCreated(newPVZ)

	handlerutils.WriteSuccess(w, http.StatusCreated, CreatePVZResponse{
		ID:               newPVZ.ID,
		RegistrationDate: newPVZ.RegistrationDate,
		City:             newPVZ.City,
	})
}
