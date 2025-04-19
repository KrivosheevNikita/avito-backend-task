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

func overrideGenerateToken(fn func(role string) (string, error)) {
	handler.GenerateTokenFunc = fn
}

func TestDummyLoginHandler_Success(t *testing.T) {
	expectedToken := "mock-token"
	overrideGenerateToken(func(role string) (string, error) {
		assert.Equal(t, "employee", role)
		return expectedToken, nil
	})

	body, _ := json.Marshal(handler.DummyLoginRequest{Role: "employee"})
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.DummyLoginHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result handler.TokenResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, result.Token)
}

func TestDummyLoginHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader([]byte("invalid-json")))
	w := httptest.NewRecorder()

	handler.DummyLoginHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestDummyLoginHandler_InvalidRole(t *testing.T) {
	overrideGenerateToken(func(role string) (string, error) {
		return "", models.ErrInvalidRole
	})

	body, _ := json.Marshal(handler.DummyLoginRequest{Role: "unknown"})
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.DummyLoginHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}
