package routes

import (
	"database/sql"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/controller"
	"github.com/BUIHUUHAIDANG/freelance_market/internal/swagger"
	"github.com/gin-gonic/gin"
)

// HealthHandler godoc
// @Summary Check API status
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Server is running"})
}

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/health", HealthHandler)
	router.POST("/auth/login", controller.LoginHandler(db))
	router.POST("/auth/register", controller.RegisterHandler(db))
	swagger.RegisterRoutes(router)
	UserRoutes(router,db)
}
