package reception

import (
	"encoding/json"
	"net/http"
	"pvz-service/internal/app/reception"
	"pvz-service/internal/models"
	"pvz-service/pkg/handlerutils"
	"pvz-service/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CreateReceptionRequest struct {
	PvzID string `json:"pvzId"`
}

type CreateReceptionResponse struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzID    uuid.UUID `json:"pvzId"`
	Status   string    `json:"status"`
}

var CreateReceptionFunc = reception.CreateReception

func CreateReceptionHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogRequest(r, "Получен запрос на создание новой приемки")

	var req CreateReceptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Неверное тело запроса", err)
		handlerutils.WriteError(w, models.ErrBadRequest)
		return
	}

	pvzID, err := uuid.Parse(req.PvzID)
	if err != nil {
		logger.Error("Неверный формат pvzId", err)
		handlerutils.WriteError(w, models.ErrPvzID)
		return
	}

	logger.Info("Начало создания новой приемки", logrus.Fields{
		"pvzId": pvzID,
	})

	rec, err := CreateReceptionFunc(pvzID)
	if err != nil {
		logger.Error("Ошибка при создании приемки", err)
		handlerutils.WriteError(w, err)
		return
	}

	logger.LogReceptionCreated(rec)

	handlerutils.WriteSuccess(w, http.StatusCreated, rec)
}
