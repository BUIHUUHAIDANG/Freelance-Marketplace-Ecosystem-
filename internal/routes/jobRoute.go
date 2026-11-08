package routes

import (
	"database/sql"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/controller"
	"github.com/BUIHUUHAIDANG/freelance_market/internal/middleware"
	"github.com/gin-gonic/gin"
)


func JobRoute(router  *gin.Engine, db *sql.DB){
	router.GET("/job",controller.GetAllJobs(db))
	router.GET("/job/:jobid",controller.GetJobById(db))

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/job",controller.AddJob(db))
		protected.PUT("/job/:jobid",controller.UpdateJobById(db))
		protected.DELETE("/job/:jobid",controller.DeleteJobById(db))
	}
}