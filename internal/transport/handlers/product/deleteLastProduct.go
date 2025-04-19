package product

import (
	"net/http"
	"pvz-service/internal/app/product"
	"pvz-service/internal/models"
	"pvz-service/pkg/handlerutils"
	"pvz-service/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var DeleteLastProductFunc = product.DeleteLastProduct

func DeleteLastProductHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogRequest(r, "Получен запрос на удаление последнего товара")

	vars := mux.Vars(r)
	pvzIDStr, ok := vars["pvzId"]
	if !ok {
		handlerutils.WriteError(w, models.ErrPvzID)
		return
	}

	pvzID, err := uuid.Parse(pvzIDStr)
	if err != nil {

		handlerutils.WriteError(w, models.ErrPvzID)
		return
	}

	logger.Info("Попытка удаления последнего товара", map[string]interface{}{
		"pvzId": pvzID,
	})

	err = DeleteLastProductFunc(pvzID)
	if err != nil {

		handlerutils.WriteError(w, err)
		return
	}

	logger.Info("Последний товар успешно удалён", map[string]interface{}{
		"pvzId": pvzID,
	})

	handlerutils.WriteSuccess(w, http.StatusOK, nil)
}
