package pvz

import (
	"net/http"
	"pvz-service/internal/app/pvz"
	"pvz-service/internal/models"
	"pvz-service/pkg/handlerutils"
	"pvz-service/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var CloseLastReceptionFunc = pvz.CloseLastReception

func CloseLastReceptionHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogRequest(r, "Получен запрос на закрытие последней приемки")

	pvzIdStr := mux.Vars(r)["pvzId"]
	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		handlerutils.WriteError(w, models.ErrPvzID)
		return
	}

	logger.Info("Начало закрытия последней приемки", map[string]interface{}{
		"pvzId": pvzId,
	})

	reception, err := CloseLastReceptionFunc(pvzId)
	if err != nil {
		handlerutils.WriteError(w, err)
		return
	}

	logger.LogReceptionClosed(reception)

	handlerutils.WriteSuccess(w, http.StatusOK, reception)
}
