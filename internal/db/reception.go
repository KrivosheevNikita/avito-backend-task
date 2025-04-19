package db

import (
	"context"
	"fmt"
	"pvz-service/internal/models"
	"pvz-service/pkg/logger"

	"github.com/google/uuid"
)

func SaveReception(reception *models.Reception) error {
	query := `
		INSERT INTO receptions (id, date_time, pvz_id, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, date_time, pvz_id, status
	`

	db := GetDBExecutor()

	err := db.QueryRow(context.Background(), query, reception.ID, reception.DateTime, reception.PvzID, reception.Status).
		Scan(&reception.ID, &reception.DateTime, &reception.PvzID, &reception.Status)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса на создание приемки", err)
		return fmt.Errorf("ошибка при выполнении запроса на создание приемки: %w", err)
	}

	return nil
}

func GetActiveReceptionByPVZ(pvzID uuid.UUID) (*models.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status
		FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		LIMIT 1
	`

	db := GetDBExecutor()

	var reception models.Reception
	err := db.QueryRow(context.Background(), query, pvzID).
		Scan(&reception.ID, &reception.DateTime, &reception.PvzID, &reception.Status)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoActiveReception
		}
		logger.Error("Ошибка при выполнении запроса на получение активной приемки", err)
		return nil, fmt.Errorf("Ошибка при выполнении запроса на получение активной приемки: %w", err)
	}

	return &reception, nil
}

func HasActiveReception(pvzID uuid.UUID) (bool, error) {
	db := GetDBExecutor()

	checkQuery := `SELECT EXISTS (SELECT 1 FROM pvz WHERE id = $1)`
	var pvzExists bool
	err := db.QueryRow(context.Background(), checkQuery, pvzID).Scan(&pvzExists)
	if err != nil {
		logger.Error("Ошибка при проверке существования ПВЗ", err)
		return false, fmt.Errorf("Ошибка при проверке существования ПВЗ: %w", err)
	}
	if !pvzExists {
		return false, models.ErrPVZNotFound
	}

	query := `SELECT EXISTS (SELECT 1 FROM receptions WHERE pvz_id = $1 AND status = 'in_progress')`
	var exists bool
	err = db.QueryRow(context.Background(), query, pvzID).Scan(&exists)
	if err != nil {
		logger.Error("Ошибка при проверке активной приемки", err)
		return false, fmt.Errorf("Ошибка при проверке активной приемки: %w", err)
	}

	return exists, nil
}

func GetReceptionsByPvzID(pvzID uuid.UUID) ([]models.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status
		FROM receptions
		WHERE pvz_id = $1
		ORDER BY date_time DESC
	`

	db := GetDBExecutor()

	rows, err := db.Query(context.Background(), query, pvzID)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса на получение приемок", err)
		return nil, fmt.Errorf("Ошибка при выполнении запроса на получение приемок: %w", err)
	}
	defer rows.Close()

	var receptions []models.Reception
	for rows.Next() {
		var reception models.Reception
		if err := rows.Scan(&reception.ID, &reception.DateTime, &reception.PvzID, &reception.Status); err != nil {
			logger.Error("Ошибка при сканировании результатов запроса", err)
			return nil, fmt.Errorf("Ошибка при сканировании результатов запроса: %w", err)
		}
		receptions = append(receptions, reception)
	}

	return receptions, nil
}

func FetchProductsByReception(receptionID uuid.UUID) ([]models.Product, error) {
	db := GetDBExecutor()

	query := `SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1 ORDER BY date_time`
	rows, err := db.Query(context.Background(), query, receptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.DateTime, &p.Type, &p.ReceptionID); err != nil {
			logger.Error("Ошибка при получении продуктов из приемки", err)
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func UpdateReception(reception *models.Reception) error {
	query := `
		UPDATE receptions
		SET status = $1
		WHERE id = $2
	`

	db := GetDBExecutor()

	_, err := db.Exec(context.Background(), query, reception.Status, reception.ID)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса на обновление приемки", err)
		return fmt.Errorf("Ошибка при выполнении запроса на обновление приемки: %w", err)
	}

	return nil
}
