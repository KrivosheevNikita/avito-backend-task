package db

import (
	"context"
	"fmt"
	"pvz-service/internal/models"
	"pvz-service/pkg/logger"
	"time"

	"github.com/google/uuid"
)

func SavePVZ(newPVZ *models.PVZ) error {
	query := `
        INSERT INTO pvz (id, city, registration_date)
        VALUES ($1, $2, $3)
    `

	exec := GetDBExecutor()

	_, err := exec.Exec(context.Background(), query, newPVZ.ID, newPVZ.City, newPVZ.RegistrationDate)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса на создание ПВЗ", err)
		return fmt.Errorf("Ошибка при выполнении запроса на создание ПВЗ: %w", err)
	}

	return nil
}

func FetchPVZ(startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error) {
	exec := GetDBExecutor()

	offset := (page - 1) * limit

	query := `
        SELECT p.id, p.city, p.registration_date 
        FROM pvz p
        JOIN (
            SELECT pvz_id, MAX(date_time) AS last_reception_date
            FROM receptions
            GROUP BY pvz_id
        ) r ON p.id = r.pvz_id
        WHERE 1=1
    `

	var args []interface{}
	paramCount := 0

	if startDate != nil {
		paramCount++
		query += fmt.Sprintf(" AND r.last_reception_date >= $%d", paramCount)
		args = append(args, *startDate)
	}
	if endDate != nil {
		paramCount++
		query += fmt.Sprintf(" AND r.last_reception_date <= $%d", paramCount)
		args = append(args, *endDate)
	}

	query += " ORDER BY r.last_reception_date DESC"

	paramCount++
	query += fmt.Sprintf(" LIMIT $%d", paramCount)
	args = append(args, limit)

	paramCount++
	query += fmt.Sprintf(" OFFSET $%d", paramCount)
	args = append(args, offset)

	rows, err := exec.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pvzs []models.PVZ
	for rows.Next() {
		var p models.PVZ
		if err := rows.Scan(&p.ID, &p.City, &p.RegistrationDate); err != nil {
			return nil, err
		}
		pvzs = append(pvzs, p)
	}

	return pvzs, nil
}

func FetchReceptionsByPVZ(pvzID uuid.UUID) ([]models.Reception, error) {
	exec := GetDBExecutor()

	query := `SELECT id, date_time, pvz_id, status FROM receptions 
			  WHERE pvz_id = $1 ORDER BY date_time DESC`
	rows, err := exec.Query(context.Background(), query, pvzID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receptions []models.Reception
	for rows.Next() {
		var r models.Reception
		if err := rows.Scan(&r.ID, &r.DateTime, &r.PvzID, &r.Status); err != nil {
			return nil, err
		}
		receptions = append(receptions, r)
	}
	return receptions, nil
}

func GetPvz() ([]models.PVZ, error) {
	query := `
		SELECT id, registration_date, city
		FROM pvz
		ORDER BY registration_date DESC
	`

	db := GetDBExecutor()

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		logger.Error("Ошибка при получении ПВЗ из БД", err)
		return nil, err
	}
	defer rows.Close()

	var pvzs []models.PVZ
	for rows.Next() {
		var pvz models.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			logger.Error("Ошибка при сканировании строки ПВЗ", err)
			return nil, err
		}
		pvzs = append(pvzs, pvz)
	}

	return pvzs, nil
}
