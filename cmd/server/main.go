package main

import (
	"log"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/config"
	"github.com/BUIHUUHAIDANG/freelance_market/internal/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// @title Freelance Marketplace API
// @version 1.0
// @description API for the Freelance Marketplace project.
// @BasePath /
// @SecurityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.ConnectWithDataBase()
	config.SetConfig()
	router := gin.Default()
	routes.RegisterRoutes(router, config.DB)
	if err := router.Run(":" + config.CONFIG.Port); err != nil {
		log.Fatal(err)
	}
}