package product

import (
	"encoding/json"
	"net/http"
	"pvz-service/internal/app/product"
	"pvz-service/internal/models"
	"pvz-service/pkg/handlerutils"
	"pvz-service/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type AddProductRequest struct {
	Type  string `json:"type"`
	PvzId string `json:"pvzId"`
}

type AddProductResponse struct {
	ID          uuid.UUID `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionID uuid.UUID `json:"receptionId"`
}

var AddProductFunc = product.AddProduct

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogRequest(r, "Получен запрос на добавление товара")

	var req AddProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Ошибка при разборе тела запроса", err)
		handlerutils.WriteError(w, models.ErrBadRequest)
		return
	}

	pvzID, err := uuid.Parse(req.PvzId)
	if err != nil {
		logger.Error("Неверный формат PvzId", err)
		handlerutils.WriteError(w, models.ErrPvzID)
		return
	}

	logger.Info("Попытка добавления товара", map[string]interface{}{
		"pvzId": pvzID,
		"type":  req.Type,
	})

	prod, err := AddProductFunc(pvzID, req.Type)
	if err != nil {
		logger.Error("Ошибка при добавлении товара", err)
		handlerutils.WriteError(w, err)
		return
	}

	logger.LogProductAdded(prod, pvzID)

	resp := AddProductResponse{
		ID:          prod.ID,
		DateTime:    prod.DateTime,
		Type:        prod.Type,
		ReceptionID: prod.ReceptionID,
	}

	handlerutils.WriteSuccess(w, http.StatusCreated, resp)
}
