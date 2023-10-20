package migrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Column struct {
	ColumnName      string  `json:"column_name"`
	OrdinalPosition int     `json:"ordinal_position"`
	IsNullable      string  `json:"is_nullable"`
	DataType        string  `json:"data_type"`
	ColumnDefault   *string `json:"column_default,omitempty"`
}

func GetColumns(tableName string) ([]Column, error) {
	tx, err := connection.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	if !TableExists(tableName) {
		return nil, fmt.Errorf("Table %s does not exist", tableName)
	}

	query := `
  SELECT
		COLUMN_NAME,
    ORDINAL_POSITION,
    IS_NULLABLE,
    DATA_TYPE,
    COLUMN_DEFAULT
  FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_NAME = $1 ;`

	rows, _ := tx.Query(context.Background(), query, tableName)
	columns, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Column])
	if err != nil {
		return nil, err
	}

	return columns, nil
}
