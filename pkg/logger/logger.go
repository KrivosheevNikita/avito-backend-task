package logger

import (
	"os"

	"net/http"

	"pvz-service/internal/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	Log.SetFormatter(&logrus.JSONFormatter{})

	Log.SetLevel(logrus.InfoLevel)

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Log.Fatal(err)
	}

	Log.SetOutput(file)
}

func Info(message string, fields logrus.Fields) {
	Log.WithFields(fields).Info(message)
}

func Error(message string, err error) {
	Log.WithField("error", err).Error(message)
}

func Warn(message string, fields logrus.Fields) {
	Log.WithFields(fields).Warn(message)
}

func LogRequest(r *http.Request, message string) {
	Log.WithFields(logrus.Fields{
		"method": r.Method,
		"url":    r.URL.Path,
	}).Info(message)
}

func LogReceptionCreated(rec *models.Reception) {
	Log.WithFields(logrus.Fields{
		"reception_id": rec.ID,
		"date_time":    rec.DateTime,
		"pvzId":        rec.PvzID,
		"status":       rec.Status,
	}).Info("Приемка успешно создана")
}

func LogPVZCreated(pvz *models.PVZ) {
	Log.WithFields(logrus.Fields{
		"id":                pvz.ID,
		"registration_date": pvz.RegistrationDate,
		"city":              pvz.City,
	}).Info("ПВЗ успешно создан")
}

func LogReceptionClosed(reception *models.Reception) {
	Log.WithFields(logrus.Fields{
		"reception_id": reception.ID,
		"pvzId":        reception.PvzID,
		"date_time":    reception.DateTime,
		"status":       reception.Status,
	}).Info("Приемка успешно закрыта")
}

func LogProductAdded(prod *models.Product, pvzID uuid.UUID) {
	Log.WithFields(logrus.Fields{
		"product_id":   prod.ID,
		"date_time":    prod.DateTime,
		"type":         prod.Type,
		"reception_id": prod.ReceptionID,
		"pvzId":        pvzID,
	}).Info("Товар успешно добавлен")
}

func LogUserRegistered(email, role string, userID uuid.UUID) {
	Log.WithFields(logrus.Fields{
		"email": email,
		"role":  role,
		"id":    userID,
	}).Info("Пользователь успешно зарегистрирован")
}
