package migrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var connection *pgxpool.Pool

// create a new connection
func OpenConnection(db_url string) {
	var err error

	if CheckConnection() {
		return
	}

	connection, err = pgxpool.New(context.Background(), db_url)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

	err = RunChecks()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database")
}

// close the connection
func CloseConnection() {
	connection.Close()
	fmt.Println("Connection closed")
}

// check if the connection is alive
func CheckConnection() bool {
	if connection == nil {
		return false
	}
	err := connection.Ping(context.Background())
	return err == nil
}

// run checks for the database
func RunChecks() error {
	if !CheckConnection() {
		return fmt.Errorf("Create connection first")
	}

	// check if the function uuid-ossp exists
	query := `
  SELECT EXISTS (
    SELECT 1 
    FROM pg_extension 
    WHERE extname = 'uuid-ossp'
  ) AS extension_exists;`

	var extensionExists bool
	err := connection.QueryRow(context.Background(), query).Scan(&extensionExists)
	if err != nil {
		return err
	}

	return nil
}
