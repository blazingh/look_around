package migrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Table struct {
	TableName string `json:"table_name"`
}

func TableExists(tableName string) bool {
	tx, err := connection.Begin(context.Background())
	if err != nil {
		return false
	}
	defer tx.Rollback(context.Background())

	query := `
    SELECT EXISTS (
      SELECT 1
      FROM information_schema.tables
      WHERE table_name = $1
    ) AS table_exists;`

	var tableExists bool
	err = tx.QueryRow(context.Background(), query, tableName).Scan(&tableExists)
	if err != nil {
		return false
	}
	return tableExists
}

func CreateTable(tableName string) error {
	tx, err := connection.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if TableExists(tableName) {
		return fmt.Errorf("Table %s already exists", tableName)
	}

	query := `
    CREATE TABLE IF NOT EXISTS $1 (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4()
    );`

	_, err = tx.Exec(context.Background(), query, tableName)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func DropTable(tableName string) error {
	tx, err := connection.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if !TableExists(tableName) {
		return fmt.Errorf("Table %s does not exist", tableName)
	}

	query := `DROP TABLE IF EXISTS $1 ;`

	_, err = tx.Exec(context.Background(), query, tableName)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func GetTables() ([]Table, error) {
	tx, err := connection.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	query := `
    SELECT table_name
    FROM information_schema.tables
    WHERE table_schema = 'public';`

	rows, _ := tx.Query(context.Background(), query)
	tables, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Table])

	if err != nil {
		return nil, err
	}

	return tables, nil
}

func AlterTable() {
	// TODO
}
