package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/model"
	"github.com/gin-gonic/gin"
)

func RequireJobManagerOrAdmin(db *sql.DB) gin.HandlerFunc {
	return requireProfileOrAdmin(db, "hire_manager")
}

func RequireFreelancerOrAdmin(db *sql.DB) gin.HandlerFunc {
	return requireProfileOrAdmin(db, "freelancer")
}

func requireProfileOrAdmin(db *sql.DB, profileTable string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDValue, exists := ctx.Get("user_id")
		userID, validUserID := userIDValue.(int64)
		if !exists || !validUserID || userID < 1 {
			ctx.JSON(http.StatusUnauthorized, model.Fail{Error: "Unauthorized"})
			ctx.Abort()
			return
		}

		isAdmin, err := isConfiguredAdmin(userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set("is_admin", isAdmin)
		if isAdmin {
			ctx.Next()
			return
		}
		if db == nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: "Database is unavailable"})
			ctx.Abort()
			return
		}

		var hasProfile bool
		query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE user_account_id = ?)", profileTable)
		if err := db.QueryRowContext(ctx.Request.Context(), query, userID).Scan(&hasProfile); err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			ctx.Abort()
			return
		}
		if !hasProfile {
			ctx.JSON(http.StatusForbidden, model.Fail{Error: "This action is not allowed for your account role"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func isConfiguredAdmin(userID int64) (bool, error) {
	isAdmin := false
	for _, value := range strings.Split(os.Getenv("ADMIN_USER_IDS"), ",") {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		adminID, err := strconv.ParseInt(value, 10, 64)
		if err != nil || adminID < 1 {
			return false, fmt.Errorf("ADMIN_USER_IDS must be a comma-separated list of positive user IDs")
		}
		isAdmin = isAdmin || adminID == userID
	}
	return isAdmin, nil
}
