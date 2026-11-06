package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/BUIHUUHAIDANG/freelance_market/internal/controller"
)


func UserRoutes(router *gin.Engine) {
	router.POST("/login", controller.LoginHandler)
	router.POST("/register", controller.RegisterHandler)

}