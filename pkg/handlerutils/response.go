package handlerutils

import (
	"encoding/json"
	"errors"
	"net/http"
	"pvz-service/internal/models"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, err error) {
	status := getStatusCode(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusInternalServerError {
		_ = json.NewEncoder(w).Encode(ErrorResponse{Message: "Внутрення ошибка сервера"})
	} else {
		_ = json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
	}
}

func WriteSuccess(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func getStatusCode(err error) int {
	switch {
	case errors.Is(err, models.ErrBadRequest),
		errors.Is(err, models.ErrInvalidCity),
		errors.Is(err, models.ErrInvalidDateRange),
		errors.Is(err, models.ErrInvalidPagination),
		errors.Is(err, models.ErrInvalidRole),
		errors.Is(err, models.ErrNoActiveReception),
		errors.Is(err, models.ErrNoProductToDelete),
		errors.Is(err, models.ErrReceptionAlreadyClosed),
		errors.Is(err, models.ErrReceptionInProgress),
		errors.Is(err, models.ErrNoOpenReception),
		errors.Is(err, models.ErrInvalidTypeProduct),
		errors.Is(err, models.ErrPVZNotFound),
		errors.Is(err, models.ErrPvzID),
		errors.Is(err, models.ErrEmailExist):
		return http.StatusBadRequest

	case errors.Is(err, models.ErrUnauthorized):
		return http.StatusUnauthorized

	case errors.Is(err, models.ErrForbidden):
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}
