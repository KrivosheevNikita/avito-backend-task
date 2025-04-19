package auth

import (
	"errors"
	"pvz-service/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTokenGen struct {
	token string
	err   error
}

func (m *mockTokenGen) GenerateToken(role string) (string, error) {
	return m.token, m.err
}

type mockHasher struct {
	hashedPassword string
	hashErr        error
	checkResult    bool
}

func (m *mockHasher) HashPassword(password string) (string, error) {
	return m.hashedPassword, m.hashErr
}

func (m *mockHasher) CheckPasswordHash(password, hash string) bool {
	return m.checkResult
}

type mockUserRepo struct {
	user         *models.User
	passwordHash string
	getErr       error
	saveErr      error
	saved        *models.User
}

func (m *mockUserRepo) CreateUser(user *models.User, password string) error {
	m.saved = user
	return m.saveErr
}

func (m *mockUserRepo) GetUserByEmail(email string) (*models.User, string, error) {
	return m.user, m.passwordHash, m.getErr
}

func setupMocks(tg tokenGenerator, h passwordHasher, r userRepository) {
	tokenGen = tg
	hasher = h
	repo = r
}

func TestGenerateToken_Success(t *testing.T) {
	mockTG := &mockTokenGen{token: "mocked-token"}
	setupMocks(mockTG, nil, nil)

	token, err := GenerateToken("employee")
	assert.NoError(t, err)
	assert.Equal(t, "mocked-token", token)
}
func TestGenerateToken_InvalidRole(t *testing.T) {
	token, err := GenerateToken("admin")
	assert.ErrorIs(t, err, models.ErrInvalidRole)
	assert.Empty(t, token)
}

func TestRegister_Success(t *testing.T) {
	mockHasher := &mockHasher{hashedPassword: "hashed-pass"}
	mockRepo := &mockUserRepo{}
	setupMocks(nil, mockHasher, mockRepo)

	user, err := Register("test@example.com", "pass", "employee")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "employee", user.Role)
	assert.Equal(t, user, mockRepo.saved)
}

func TestRegister_HashError(t *testing.T) {
	mockHasher := &mockHasher{hashErr: errors.New("hash error")}
	setupMocks(nil, mockHasher, nil)

	user, err := Register("a@a.com", "pass", "employee")
	assert.Nil(t, user)
	assert.EqualError(t, err, "hash error")
}

func TestRegister_SaveError(t *testing.T) {
	mockHasher := &mockHasher{hashedPassword: "pass"}
	mockRepo := &mockUserRepo{saveErr: errors.New("save fail")}
	setupMocks(nil, mockHasher, mockRepo)

	user, err := Register("x@x.com", "123", "moderator")
	assert.Nil(t, user)
	assert.EqualError(t, err, "save fail")
}

func TestRegister_InvalidRole(t *testing.T) {
	user, err := Register("x@x.com", "123", "admin")
	assert.Nil(t, user)
	assert.ErrorIs(t, err, models.ErrInvalidRole)
}

func TestLogin_Success(t *testing.T) {
	mockUser := &models.User{Email: "test@example.com", Role: "employee"}
	mockRepo := &mockUserRepo{
		user:         mockUser,
		passwordHash: "hashed-pass",
	}
	mockHasher := &mockHasher{checkResult: true}
	mockTG := &mockTokenGen{token: "jwt-token"}

	setupMocks(mockTG, mockHasher, mockRepo)

	token, err := Login("test@example.com", "password")
	assert.NoError(t, err)
	assert.Equal(t, "jwt-token", token)
}

func TestLogin_GetUserError(t *testing.T) {
	mockRepo := &mockUserRepo{getErr: errors.New("not found")}
	setupMocks(nil, nil, mockRepo)

	token, err := Login("no@user.com", "123")
	assert.ErrorIs(t, err, models.ErrUnauthorized)
	assert.Empty(t, token)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockUser := &models.User{Email: "user@example.com", Role: "employee"}
	mockRepo := &mockUserRepo{
		user:         mockUser,
		passwordHash: "hashed-pass",
	}
	mockHasher := &mockHasher{checkResult: false}
	setupMocks(nil, mockHasher, mockRepo)

	token, err := Login("user@example.com", "wrong")
	assert.ErrorIs(t, err, models.ErrUnauthorized)
	assert.Empty(t, token)
}
