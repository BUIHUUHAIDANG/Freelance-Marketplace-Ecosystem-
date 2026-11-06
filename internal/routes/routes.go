package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is running",
		})
	})
	router.POST("/login", controller.LoginHandler)
	router.POST("/register", controller.RegisterHandler)
	UserRoutes(router)
}
