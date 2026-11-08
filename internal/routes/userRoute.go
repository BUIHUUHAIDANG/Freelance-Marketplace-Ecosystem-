package routes

import (
	"database/sql"
	"github.com/BUIHUUHAIDANG/freelance_market/internal/controller"
	"github.com/BUIHUUHAIDANG/freelance_market/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, db *sql.DB) {
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/user/:id", controller.GetUserProfile(db))
		protected.PUT("/user/:id", controller.UpdateUserProfile(db))
		protected.DELETE("/user/:id", controller.DeleteUserProfile(db))
		protected.GET("/user/:id/freelancer", controller.GetUserFreelancerProfile(db))
		protected.GET("/user/:id/hirer", controller.GetUserHirerProfile(db))
	}
}
