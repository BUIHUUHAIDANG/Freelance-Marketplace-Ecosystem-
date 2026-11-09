package routes

import (
	"database/sql"
	"github.com/BUIHUUHAIDANG/freelance_market/internal/controller"
	"github.com/BUIHUUHAIDANG/freelance_market/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, db *sql.DB) {
	protected := router.Group("/")
	protected.Use(middleware.JwtAuthMiddleware())
	{
		protected.GET("/user/me", controller.GetUserProfile(db))
		protected.PUT("/user/me", controller.UpdateUserProfile(db))
		protected.DELETE("/user/me", controller.DeleteUserProfile(db))
		protected.GET("/user/me/freelancer", controller.GetUserFreelancerProfile(db))
		protected.GET("/user/me/hirer", controller.GetUserHirerProfile(db))
	}
}
