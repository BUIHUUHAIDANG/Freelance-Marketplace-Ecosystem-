package main

import (
	_ "fmt"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/config"
	"github.com/gin-gonic/gin"
)



func main(){
	cfg,_ := config.LoadConfig();
	router := gin.Default()
  	router.GET("/ping", func(c *gin.Context) {
  	  c.JSON(200, gin.H{
  	    "message": "pong",
  	  })
  	})
  	router.Run(":"+cfg.Port)
}