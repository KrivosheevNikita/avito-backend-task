package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type mockDBExecutor struct {
	mock.Mock
}

func (m *mockDBExecutor) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	argsArray := m.Called(ctx, sql, args)
	return argsArray.Get(0).(pgconn.CommandTag), argsArray.Error(1)
}

func (m *mockDBExecutor) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	argsArray := m.Called(ctx, sql, args)

	rowsInterface := argsArray.Get(0)
	var rows pgx.Rows
	if rowsInterface != nil {
		rows = rowsInterface.(pgx.Rows)
	}

	return rows, argsArray.Error(1)
}

func (m *mockDBExecutor) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	arguments := m.Called(append([]interface{}{ctx, sql}, args...)...)
	return arguments.Get(0).(pgx.Row)
}

type mockDBRow struct {
	scanFn func(dest ...interface{}) error
}

func (m mockDBRow) Scan(dest ...interface{}) error {
	return m.scanFn(dest...)
}
