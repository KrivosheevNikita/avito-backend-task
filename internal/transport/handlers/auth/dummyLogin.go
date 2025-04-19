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

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

var GenerateTokenFunc = auth.GenerateToken

func DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogRequest(r, "Получен запрос на Dummy-авторизацию")

	var req DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Ошибка при разборе тела запроса", err)
		handlerutils.WriteError(w, models.ErrBadRequest)
		return
	}

	logger.Info("Попытка генерации токена", logrus.Fields{
		"role": req.Role,
	})

	token, err := GenerateTokenFunc(req.Role)
	if err != nil {
		logger.Error("Ошибка при генерации токена", err)
		handlerutils.WriteError(w, err)
		return
	}

	logger.Info("Токен успешно сгенерирован", logrus.Fields{
		"role": req.Role,
	})

	handlerutils.WriteSuccess(w, http.StatusOK, TokenResponse{Token: token})
}
