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

func ColumnExists(tableName string, columnName string) bool {
	if !TableExists(tableName) {
		return false
	}
	query := `
    SELECT EXISTS (
      SELECT 1
      FROM information_schema.columns
      WHERE table_name = $1 AND column_name = $2
    ) AS column_exists;`

	var columnExists bool
	err := connection.QueryRow(context.Background(), query, tableName, columnName).Scan(&columnExists)
	if err != nil {
		return false
	}
	return columnExists
}

func CreateColumn(tableName string, column Column) error {
	tx, err := connection.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if !ColumnExists(tableName, column.ColumnName) {
		return fmt.Errorf("Column %s does not exist", column.ColumnName)
	}

	query := `
    ALTER TABLE $1
    ADD COLUMN $2;`

	_, err = tx.Exec(context.Background(), query, tableName, column.ColumnName)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func DropColumn(tableName string, columnName string) error {
	tx, err := connection.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if !ColumnExists(tableName, columnName) {
		return fmt.Errorf("Column %s does not exist", columnName)
	}

	query := `
    ALTER TABLE $1
    DROP COLUMN $2;`

	_, err = tx.Exec(context.Background(), query, tableName, columnName)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
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

func AlterColumn() {
	// TODO
}
