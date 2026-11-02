package routes

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	mysql "github.com/go-sql-driver/mysql"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	api := router.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) {
		if err := db.PingContext(c.Request.Context()); err != nil {
			writeDatabaseError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	registerJobRoutes(api, db)
}

func writeDatabaseError(c *gin.Context, err error) {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1452:
			c.JSON(http.StatusBadRequest, gin.H{"error": "one or more referenced records do not exist"})
			return
		case 1451, 1062:
			c.JSON(http.StatusConflict, gin.H{"error": "the operation conflicts with existing records"})
			return
		}
	}

	log.Printf("database request failed: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "database request failed"})
}

func writeNotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, gin.H{"error": resource + " not found"})
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return 0, false
	}
	return id, true
}
