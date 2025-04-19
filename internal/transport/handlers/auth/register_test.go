package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pvz-service/internal/models"
	handler "pvz-service/internal/transport/handlers/auth"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func overrideRegister(fn func(email, password, role string) (*models.User, error)) {
	handler.RegisterFunc = fn
}

func TestRegisterHandler_Success(t *testing.T) {
	expectedUser := &models.User{
		ID:    uuid.New(),
		Email: "email",
		Role:  "employee",
	}

	overrideRegister(func(email, password, role string) (*models.User, error) {
		assert.Equal(t, "email", email)
		assert.Equal(t, "password", password)
		assert.Equal(t, "employee", role)
		return expectedUser, nil
	})

	body, _ := json.Marshal(handler.RegisterRequest{
		Email:    "email",
		Password: "password",
		Role:     "employee",
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.RegisterHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var user models.User
	err := json.NewDecoder(resp.Body).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Role, user.Role)
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("{invalid_json")))
	w := httptest.NewRecorder()

	handler.RegisterHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestRegisterHandler_BadRequest(t *testing.T) {
	overrideRegister(func(email, password, role string) (*models.User, error) {
		return nil, models.ErrInvalidRole
	})

	body, _ := json.Marshal(handler.RegisterRequest{
		Email:    "email",
		Password: "pass",
		Role:     "r",
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.RegisterHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}
