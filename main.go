package main

import (
	"CRUD/model/db/sql"
	"CRUD/server"
	"fmt"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found")
		return
	}

	db, err := sql.NewPostgresDb()
	if err != nil {
		fmt.Println(err)
		return
	}

	server := server.NewServer(db)
	server.Start()
}
