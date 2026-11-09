package middleware

import (
	"net/http"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/model"
	"github.com/BUIHUUHAIDANG/freelance_market/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement your authentication logic here
		// For example, check for a valid token in the request header
		// If the token is valid, proceed to the next handler
		c.Next()
	}
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := utils.ExtractTokenUserID(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, model.Fail{
				Error: "Unauthorized",
			})
			ctx.Abort()
			return
		}
		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
