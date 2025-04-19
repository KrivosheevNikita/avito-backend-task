package pvz

import (
	"net/http"
	"pvz-service/internal/app/pvz"
	"pvz-service/internal/models"
	"pvz-service/pkg/handlerutils"
	"pvz-service/pkg/logger"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type GetPVZRequest struct {
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	Limit     int
}

var GetPVZListFunc = pvz.GetPVZList

func GetPVZHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogRequest(r, "Получен запрос на получение списка ПВЗ")

	var request GetPVZRequest
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if startDateStr != "" {
		t, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			handlerutils.WriteError(w, models.ErrInvalidDateRange)
			return
		}
		request.StartDate = &t
	}

	if endDateStr != "" {
		t, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			handlerutils.WriteError(w, models.ErrInvalidDateRange)
			return
		}
		request.EndDate = &t
	}

	var limit, page int
	var err error

	if limitStr == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 30 {
			handlerutils.WriteError(w, models.ErrInvalidPagination)
			return
		}
	}
	request.Limit = limit

	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			handlerutils.WriteError(w, models.ErrInvalidPagination)
			return
		}
	}
	request.Page = page

	logger.Info("Начало получения списка ПВЗ", logrus.Fields{
		"startDate": request.StartDate,
		"endDate":   request.EndDate,
		"page":      request.Page,
		"limit":     request.Limit,
	})

	result, err := GetPVZListFunc(request.StartDate, request.EndDate, request.Page, request.Limit)
	if err != nil {
		handlerutils.WriteError(w, err)
		return
	}

	logger.Info("Список ПВЗ успешно получен", logrus.Fields{
		"кол-во_ПВЗ": len(result),
	})

	handlerutils.WriteSuccess(w, http.StatusOK, result)
}
