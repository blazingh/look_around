package main

import (
	"fmt"

	"github.com/blazingh/look_around/pkg/migrations"
)

func main() {
	fmt.Println("Hello World!")
	migrations.OpenConnection("postgres://postgres:password@localhost:5432/postgres")
	defer migrations.CloseConnection()
}
