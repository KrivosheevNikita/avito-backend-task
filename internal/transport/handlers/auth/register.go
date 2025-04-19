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

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var RegisterFunc = auth.Register

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogRequest(r, "Получен запрос на регистрацию")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Ошибка при разборе тела запроса на регистрацию", err)
		handlerutils.WriteError(w, models.ErrBadRequest)
		return
	}

	logger.Info("Обработка запроса на регистрацию", logrus.Fields{
		"email": req.Email,
		"role":  req.Role,
	})

	user, err := RegisterFunc(req.Email, req.Password, req.Role)
	if err != nil {
		logger.Error("Не удалось зарегистрироваться", err)
		handlerutils.WriteError(w, err)
		return
	}

	logger.LogUserRegistered(req.Email, req.Role, user.ID)

	handlerutils.WriteSuccess(w, http.StatusCreated, user)
}
