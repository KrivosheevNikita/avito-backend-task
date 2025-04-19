package auth

import (
	"encoding/json"
	"net/http"
	"pvz-service/internal/app/auth"
	"pvz-service/internal/models"
	"pvz-service/pkg/handlerutils"
	"pvz-service/pkg/logger"

	"github.com/sirupsen/logrus"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var LoginFunc = auth.Login

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogRequest(r, "Получен запрос на авторизацию")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Ошибка при разборе тела запроса на вход", err)
		handlerutils.WriteError(w, models.ErrBadRequest)
		return
	}

	logger.Info("Обработка запроса на вход", logrus.Fields{
		"email": req.Email,
	})

	token, err := LoginFunc(req.Email, req.Password)
	if err != nil {
		logger.Error("Неверные учетные данные", err)
		handlerutils.WriteError(w, err)
		return
	}

	logger.Info("Успешный вход", logrus.Fields{
		"email": req.Email,
	})

	handlerutils.WriteSuccess(w, http.StatusOK, TokenResponse{Token: token})
}
