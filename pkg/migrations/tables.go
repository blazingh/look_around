package migrations

import (
	"context"
	"fmt"
)

func CreateTable(tableName string) error {
	tx, err := connection.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// check if the table already exists
	query := `
    SELECT EXISTS (
      SELECT 1
      FROM information_schema.tables
      WHERE table_name = '` + tableName + `'
    ) AS table_exists;`

	var tableExists bool
	err = tx.QueryRow(context.Background(), query).Scan(&tableExists)
	if err != nil {
		return err
	}
	if tableExists {
		return fmt.Errorf("Table %s already exists", tableName)
	}

	query = `
    CREATE TABLE IF NOT EXISTS ` + tableName + ` (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4()
    );`

	_, err = tx.Exec(context.Background(), query)
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

	// check if the table already exists
	query := `
    SELECT EXISTS (
      SELECT 1
      FROM information_schema.tables
      WHERE table_name = '` + tableName + `'
    ) AS table_exists;`

	var tableExists bool
	err = tx.QueryRow(context.Background(), query).Scan(&tableExists)
	if err != nil {
		return err
	}
	if !tableExists {
		return fmt.Errorf("Table %s does not exist", tableName)
	}

	query = `
    DROP TABLE IF EXISTS ` + tableName + `;`

	_, err = tx.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func AlterTable() {
	// TODO
}
