package auth

import (
	"pvz-service/internal/db"
	"pvz-service/internal/models"
	"pvz-service/pkg/utils"

	"github.com/google/uuid"
)

type tokenGenerator interface {
	GenerateToken(role string) (string, error)
}

type passwordHasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type userRepository interface {
	CreateUser(user *models.User, password string) error
	GetUserByEmail(email string) (*models.User, string, error)
}

type defaultTokenGenerator struct{}

func (defaultTokenGenerator) GenerateToken(role string) (string, error) {
	return utils.GenerateToken(role)
}

type defaultPasswordHasher struct{}

func (defaultPasswordHasher) HashPassword(password string) (string, error) {
	return utils.HashPassword(password)
}

func (defaultPasswordHasher) CheckPasswordHash(password, hash string) bool {
	return utils.CheckPasswordHash(password, hash)
}

type defaultUserRepository struct{}

func (defaultUserRepository) CreateUser(user *models.User, password string) error {
	return db.CreateUser(user, password)
}

func (defaultUserRepository) GetUserByEmail(email string) (*models.User, string, error) {
	return db.GetUserByEmail(email)
}

var (
	tokenGen tokenGenerator = defaultTokenGenerator{}
	hasher   passwordHasher = defaultPasswordHasher{}
	repo     userRepository = defaultUserRepository{}
)

func GenerateToken(role string) (string, error) {
	if role != "employee" && role != "moderator" {
		return "", models.ErrInvalidRole
	}
	return tokenGen.GenerateToken(role)
}

func Register(email, password, role string) (*models.User, error) {
	if role != "employee" && role != "moderator" {
		return nil, models.ErrInvalidRole
	}

	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:    uuid.New(),
		Email: email,
		Role:  role,
	}

	if err := repo.CreateUser(user, hashedPassword); err != nil {
		return nil, err
	}

	return user, nil
}

func Login(email, password string) (string, error) {
	user, hash_pass, err := repo.GetUserByEmail(email)
	if err != nil {
		return "", models.ErrUnauthorized
	}

	if !hasher.CheckPasswordHash(password, hash_pass) {
		return "", models.ErrUnauthorized
	}

	return tokenGen.GenerateToken(user.Role)
}
