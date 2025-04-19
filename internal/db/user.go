package db

import (
	"context"
	"pvz-service/internal/models"
)

func CreateUser(user *models.User, password string) error {
	db := GetDBExecutor()

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := db.QueryRow(context.Background(), checkQuery, user.Email).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return models.ErrEmailExist
	}

	query := `
		INSERT INTO users (id, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
	`

	_, err = db.Exec(context.Background(), query, user.ID, user.Email, password, user.Role)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (*models.User, string, error) {
	db := GetDBExecutor()

	query := `
		SELECT id, email, role, password_hash
		FROM users
		WHERE email = $1
	`

	row := db.QueryRow(context.Background(), query, email)

	var user models.User
	var password string
	err := row.Scan(&user.ID, &user.Email, &user.Role, &password)

	if err != nil {
		return nil, "", err
	}

	return &user, password, nil
}
