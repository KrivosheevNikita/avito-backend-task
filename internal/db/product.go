package db

import (
	"context"
	"fmt"
	"pvz-service/internal/models"
	"pvz-service/pkg/logger"

	"github.com/google/uuid"
)

func AddProductToReception(product *models.Product) error {
	query := `
        INSERT INTO products (id, date_time, type, reception_id)
        VALUES ($1, $2, $3, $4)
    `

	exec := GetDBExecutor()
	_, err := exec.Exec(
		context.Background(),
		query,
		product.ID,
		product.DateTime,
		product.Type,
		product.ReceptionID,
	)

	if err != nil {
		logger.Error("Ошибка при выполнении запроса на добавление товара", err)
		return fmt.Errorf("Ошибка при выполнении запроса на добавление товара: %w", err)
	}

	return nil
}

func DeleteLastProductByReceptionID(receptionID uuid.UUID) error {
	query := `
        DELETE FROM products
        WHERE id = (
            SELECT id FROM products
            WHERE reception_id = $1
            ORDER BY date_time DESC
            LIMIT 1
        )`

	exec := GetDBExecutor()
	cmd, err := exec.Exec(context.Background(), query, receptionID)
	if err != nil {
		logger.Error("Ошибка при удалении последнего товара", err)
		return fmt.Errorf("Ошибка при удалении последнего товара: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		logger.Error("Нет товаров для удаления", models.ErrNoProductToDelete)
		return models.ErrNoProductToDelete
	}
	return nil
}
