package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"pvz-service/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbPool     *pgxpool.Pool
	dbExecutor DBExecutor
	once       sync.Once
)

type DBExecutor interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	var err error

	once.Do(func() {
		connString := cfg.DB.ConnectionString()
		dbPool, err = pgxpool.New(context.Background(), connString)
		if err != nil {
			err = fmt.Errorf("Не удается подключиться к бд: %v", err)
			return
		}

		if pingErr := dbPool.Ping(context.Background()); pingErr != nil {
			err = fmt.Errorf("Не удается проверить соединение: %v", pingErr)
			return
		}

		log.Println("Успешное подключение к БД")

		dbExecutor = dbPool
	})

	return dbPool, err
}

func SetDBExecutor(exec DBExecutor) {
	dbExecutor = exec
}

func GetDBExecutor() DBExecutor {
	if dbExecutor == nil {
		panic("DBExecutor не инициализирован")
	}
	return dbExecutor
}
