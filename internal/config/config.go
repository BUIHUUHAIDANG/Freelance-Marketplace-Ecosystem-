package config

import (
	"log"
	"os"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/model"
	"github.com/joho/godotenv"
)

var CONFIG *model.Config

func SetConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Can not load file .env")
	}
	CONFIG = &model.Config{}
	CONFIG.Port = os.Getenv("PORT")
	if CONFIG.Port == "" {
		log.Fatal("PORT is required")
	}
}