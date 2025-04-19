package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pvz-service/internal/models"
	handler "pvz-service/internal/transport/handlers/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func overrideLogin(fn func(email, password string) (string, error)) {
	handler.LoginFunc = fn
}

func TestLoginHandler_Success(t *testing.T) {
	expectedToken := "token"

	overrideLogin(func(email, password string) (string, error) {
		assert.Equal(t, "email", email)
		assert.Equal(t, "password", password)
		return expectedToken, nil
	})

	reqBody := handler.LoginRequest{
		Email:    "email",
		Password: "password",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.LoginHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response handler.TokenResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, response.Token)
}

func TestLoginHandler_Invalid(t *testing.T) {
	overrideLogin(func(email, password string) (string, error) {
		return "", models.ErrUnauthorized
	})

	reqBody := handler.LoginRequest{
		Email:    "email",
		Password: "wrong",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.LoginHandler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("{invalid_json")))
	w := httptest.NewRecorder()

	handler.LoginHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}
