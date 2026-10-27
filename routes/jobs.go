package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type job struct {
	ID            int64    `json:"id"`
	HireManagerID int64    `json:"hire_manager_id" binding:"required"`
	DurationID    int64    `json:"duration_id" binding:"required"`
	ComplexityID  int64    `json:"complexity_id" binding:"required"`
	Description   *string  `json:"description"`
	MainSkillID   int64    `json:"main_skill_id" binding:"required"`
	PaymentTypeID int64    `json:"payment_type_id" binding:"required"`
	PaymentAmount *float64 `json:"payment_amount"`
	Title         string   `json:"title" binding:"required"`
}

func registerJobRoutes(api *gin.RouterGroup, db *sql.DB) {
	api.GET("/jobs", func(c *gin.Context) {
		rows, err := db.QueryContext(c.Request.Context(), `
			SELECT id, hire_manager_id, duration_id, complexity_id, description,
			       main_skill_id, payment_type_id, payment_amount, title
			FROM job ORDER BY id DESC`)
		if err != nil {
			writeDatabaseError(c, err)
			return
		}
		defer rows.Close()

		jobs := make([]job, 0)
		for rows.Next() {
			var item job
			if err := rows.Scan(
				&item.ID, &item.HireManagerID, &item.DurationID, &item.ComplexityID,
				&item.Description, &item.MainSkillID, &item.PaymentTypeID,
				&item.PaymentAmount, &item.Title,
			); err != nil {
				writeDatabaseError(c, err)
				return
			}
			jobs = append(jobs, item)
		}
		if err := rows.Err(); err != nil {
			writeDatabaseError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": jobs})
	})

	api.GET("/jobs/:id", func(c *gin.Context) {
		id, ok := parseID(c)
		if !ok {
			return
		}
		item, err := getJob(c, db, id)
		if errors.Is(err, sql.ErrNoRows) {
			writeNotFound(c, "job")
			return
		}
		if err != nil {
			writeDatabaseError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": item})
	})

	api.POST("/jobs", func(c *gin.Context) {
		var item job
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
			return
		}
		if err := validateJob(item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := db.ExecContext(c.Request.Context(), `
			INSERT INTO job
				(hire_manager_id, duration_id, complexity_id, description,
				 main_skill_id, payment_type_id, payment_amount, title)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			item.HireManagerID, item.DurationID, item.ComplexityID, item.Description,
			item.MainSkillID, item.PaymentTypeID, item.PaymentAmount, strings.TrimSpace(item.Title))
		if err != nil {
			writeDatabaseError(c, err)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			writeDatabaseError(c, err)
			return
		}
		item, err = getJob(c, db, id)
		if err != nil {
			writeDatabaseError(c, err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data": item})
	})

	api.PUT("/jobs/:id", func(c *gin.Context) {
		id, ok := parseID(c)
		if !ok {
			return
		}
		var item job
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
			return
		}
		if err := validateJob(item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.ExecContext(c.Request.Context(), `
			UPDATE job
			SET hire_manager_id = ?, duration_id = ?, complexity_id = ?, description = ?,
			    main_skill_id = ?, payment_type_id = ?, payment_amount = ?, title = ?
			WHERE id = ?`,
			item.HireManagerID, item.DurationID, item.ComplexityID, item.Description,
			item.MainSkillID, item.PaymentTypeID, item.PaymentAmount, strings.TrimSpace(item.Title), id)
		if err != nil {
			writeDatabaseError(c, err)
			return
		}
		item, err = getJob(c, db, id)
		if errors.Is(err, sql.ErrNoRows) {
			writeNotFound(c, "job")
			return
		}
		if err != nil {
			writeDatabaseError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": item})
	})

	api.DELETE("/jobs/:id", func(c *gin.Context) {
		id, ok := parseID(c)
		if !ok {
			return
		}
		result, err := db.ExecContext(c.Request.Context(), "DELETE FROM job WHERE id = ?", id)
		if err != nil {
			writeDatabaseError(c, err)
			return
		}
		affected, err := result.RowsAffected()
		if err != nil {
			writeDatabaseError(c, err)
			return
		}
		if affected == 0 {
			writeNotFound(c, "job")
			return
		}
		c.Status(http.StatusNoContent)
	})
}

func getJob(c *gin.Context, db *sql.DB, id int64) (job, error) {
	var item job
	err := db.QueryRowContext(c.Request.Context(), `
		SELECT id, hire_manager_id, duration_id, complexity_id, description,
		       main_skill_id, payment_type_id, payment_amount, title
		FROM job WHERE id = ?`, id).Scan(
		&item.ID, &item.HireManagerID, &item.DurationID, &item.ComplexityID,
		&item.Description, &item.MainSkillID, &item.PaymentTypeID,
		&item.PaymentAmount, &item.Title,
	)
	return item, err
}

func validateJob(item job) error {
	if strings.TrimSpace(item.Title) == "" || len(item.Title) > 255 {
		return errors.New("title is required and must be at most 255 characters")
	}
	if item.HireManagerID < 1 || item.DurationID < 1 || item.ComplexityID < 1 ||
		item.MainSkillID < 1 || item.PaymentTypeID < 1 {
		return errors.New("all related IDs must be positive integers")
	}
	if item.PaymentAmount != nil && (*item.PaymentAmount <= 0 || *item.PaymentAmount > 99999999.99) {
		return errors.New("payment_amount must be greater than zero and at most 99999999.99")
	}
	return nil
}
