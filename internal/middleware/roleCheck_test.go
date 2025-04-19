package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func middleware(authHeader string, mockFunc func(string) (jwt.MapClaims, error), requiredRole string) *httptest.ResponseRecorder {
	original := parseToken
	parseToken = mockFunc
	defer func() { parseToken = original }()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	rr := httptest.NewRecorder()

	handler := RoleCheck(requiredRole)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	return rr
}

func TestRoleCheck_MissingAuthorizationHeader(t *testing.T) {
	rr := middleware("", func(token string) (jwt.MapClaims, error) {
		return nil, nil
	}, "moderator")

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestRoleCheck_InvalidToken(t *testing.T) {
	rr := middleware("Bearer invalid_token", func(token string) (jwt.MapClaims, error) {
		return nil, errors.New("Неправильный токен")
	}, "moderator")

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestRoleCheck_MissingRoleClaim(t *testing.T) {
	rr := middleware("Bearer token", func(token string) (jwt.MapClaims, error) {
		return jwt.MapClaims{}, nil
	}, "moderator")

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestRoleCheck_WrongRole(t *testing.T) {
	rr := middleware("Bearer valid_token", func(token string) (jwt.MapClaims, error) {
		return jwt.MapClaims{"role": "user"}, nil
	}, "moderator")

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestRoleCheck_CorrectRole(t *testing.T) {
	rr := middleware("Bearer valid_token", func(token string) (jwt.MapClaims, error) {
		return jwt.MapClaims{"role": "moderator"}, nil
	}, "moderator")

	assert.Equal(t, http.StatusOK, rr.Code)
}
