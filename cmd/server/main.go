package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/BUIHUUHAIDANG/freelance_market/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("load .env: %v", err)
	}

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("MYSQL_DSN is required")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	routes.RegisterRoutes(router, db)

	if err := router.Run(":8085"); err != nil {
		log.Fatal(err)
	}
}
