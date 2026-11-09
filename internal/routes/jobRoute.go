package routes

import (
	"database/sql"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/controller"
	"github.com/BUIHUUHAIDANG/freelance_market/internal/middleware"
	"github.com/gin-gonic/gin"
)

func JobRoute(router *gin.Engine, db *sql.DB) {
	router.GET("/job", controller.GetAllJobs(db))
	router.GET("/job/:jobid", controller.GetJobById(db))

	freelancerApplications := router.Group("/")
	freelancerApplications.Use(middleware.JwtAuthMiddleware())
	freelancerApplications.POST("/job/:jobid/apply",
		middleware.RequireFreelancerOrAdmin(db),
		controller.ApplyJob(db),
	)

	jobManagers := router.Group("/")
	jobManagers.Use(middleware.JwtAuthMiddleware(), middleware.RequireJobManagerOrAdmin(db))
	jobManagers.POST("/job", controller.AddJob(db))
	jobManagers.PUT("/job/:jobid", controller.UpdateJobById(db))
	jobManagers.DELETE("/job/:jobid", controller.DeleteJobById(db))
}
