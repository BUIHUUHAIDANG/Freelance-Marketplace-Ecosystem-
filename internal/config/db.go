package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectWithDataBase() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Can not load file .env")
	}
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("MYSQL_DSN is required")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatal(err)
	}
	DB = db
}
