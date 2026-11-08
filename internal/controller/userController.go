package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/model"
	"github.com/gin-gonic/gin"
)

// GetUserProfile godoc
// @Summary Get a user account profile
// @Description Returns basic account information and whether freelancer or hirer profiles exist.
// @Tags user account
// @Produce json
// @Param id path int true "User account ID"
// @Success 200 {object} model.UserAccountProfileRespond
// @Failure 400 {object} model.Fail
// @Failure 404 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /user/{id} [get]
func GetUserProfile(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userAccountProfile model.UserAccountProfileRespond
		userID, ok := parseUserID(ctx)
		if !ok {
			return
		}
		query := `SELECT 
				    ua.email AS email,
				    ua.first_name AS first_name,
				    ua.last_name AS last_name,
				    IF(f.id IS NOT NULL, TRUE, FALSE) AS has_freelancer_profile,
				    IF(hm.id IS NOT NULL, TRUE, FALSE) AS has_hirer_profile
				FROM user_account ua
				LEFT JOIN freelancer f ON ua.id = f.user_account_id
				LEFT JOIN hire_manager hm ON ua.id = hm.user_account_id
				WHERE ua.id = ?`
		err := db.QueryRowContext(ctx.Request.Context(), query, userID).
			Scan(&userAccountProfile.Email,
				&userAccountProfile.FirstName,
				&userAccountProfile.LastName,
				&userAccountProfile.HasFreelancerProfile,
				&userAccountProfile.HasHirerProfile,
			)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(404, model.Fail{
					Error: "User Profile Not Found",
				})
				return
			}
			ctx.JSON(500, model.Fail{
				Error: err.Error(),
			})
			return
		}
		ctx.JSON(200, userAccountProfile)
	}
}

// UpdateUserProfile godoc
// @Summary Update a user account profile
// @Tags user account
// @Accept json
// @Produce json
// @Param id path int true "User account ID"
// @Param request body model.UpdateUserProfileRequest true "Profile fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} model.Fail
// @Failure 404 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /user/{id} [put]
func UpdateUserProfile(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := parseUserID(ctx)
		if !ok {
			return
		}
		var input model.UpdateUserProfileRequest

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var queryParts []string
		var args []interface{}

		if input.FirstName != nil {
			queryParts = append(queryParts, "first_name = ?")
			args = append(args, *input.FirstName)
		}
		if input.LastName != nil {
			queryParts = append(queryParts, "last_name = ?")
			args = append(args, *input.LastName)
		}

		if len(queryParts) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		args = append(args, id)

		query := fmt.Sprintf("UPDATE user_account SET %s WHERE id = ?", strings.Join(queryParts, ", "))

		result, err := db.ExecContext(ctx.Request.Context(), query, args...)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{
				Error: err.Error(),
			})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{
				Error: err.Error(),
			})
			return
		}

		if rowsAffected == 0 {
			var exists bool
			if err := db.QueryRowContext(ctx.Request.Context(), "SELECT EXISTS(SELECT 1 FROM user_account WHERE id = ?)", id).Scan(&exists); err != nil {
				ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
				return
			}
			if exists {
				ctx.JSON(http.StatusOK, gin.H{"message": "User Profile Updated Successfully"})
				return
			}
			ctx.JSON(http.StatusNotFound, model.Fail{
				Error: "User Profile Not Found",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User Profile Updated Successfully"})
	}
}

// DeleteUserProfile godoc
// @Summary Delete a user account
// @Tags user account
// @Produce json
// @Param id path int true "User account ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} model.Fail
// @Failure 404 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /user/{id} [delete]
func DeleteUserProfile(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := parseUserID(ctx)
		if !ok {
			return
		}
		query := `DELETE FROM user_account
					WHERE id = ?
		`
		result, err := db.ExecContext(ctx.Request.Context(), query, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{
				Error: err.Error(),
			})
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{
				Error: err.Error(),
			})
			return
		}
		if rowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, model.Fail{
				Error: "User Profile Not Found",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "User Profile Deleted Successfully"})
	}
}

// GetUserFreelancerProfile godoc
// @Summary Get a user's freelancer profile
// @Tags user account
// @Produce json
// @Param id path int true "User account ID"
// @Success 200 {object} model.FreelancerProfileRespond
// @Failure 400 {object} model.Fail
// @Failure 404 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /user/{id}/freelancer [get]
func GetUserFreelancerProfile(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := parseUserID(ctx)
		if !ok {
			return
		}

		var profile model.FreelancerProfileRespond
		var registrationDate sql.NullTime
		err := db.QueryRowContext(ctx.Request.Context(), `
			SELECT id, user_account_id, registration_date, location, overview
			FROM freelancer
			WHERE user_account_id = ?`, userID).Scan(
			&profile.ID,
			&profile.UserAccountID,
			&registrationDate,
			&profile.Location,
			&profile.Overview,
		)
		if err != nil {
			writeProfileQueryError(ctx, err, "Freelancer Profile Not Found")
			return
		}
		if registrationDate.Valid {
			date := registrationDate.Time.Format("2006-01-02")
			profile.RegistrationDate = &date
		}
		ctx.JSON(http.StatusOK, profile)
	}
}

// GetUserHirerProfile godoc
// @Summary Get a user's hirer profile
// @Tags user account
// @Produce json
// @Param id path int true "User account ID"
// @Success 200 {object} model.HirerProfileRespond
// @Failure 400 {object} model.Fail
// @Failure 404 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /user/{id}/hirer [get]
func GetUserHirerProfile(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := parseUserID(ctx)
		if !ok {
			return
		}

		var profile model.HirerProfileRespond
		var registrationDate sql.NullTime
		err := db.QueryRowContext(ctx.Request.Context(), `
			SELECT hm.id, hm.user_account_id, hm.registration_date, hm.location,
			       c.id, c.company_name, c.company_location
			FROM hire_manager hm
			JOIN company c ON c.id = hm.company_id
			WHERE hm.user_account_id = ?`, userID).Scan(
			&profile.ID,
			&profile.UserAccountID,
			&registrationDate,
			&profile.Location,
			&profile.CompanyID,
			&profile.CompanyName,
			&profile.CompanyLocation,
		)
		if err != nil {
			writeProfileQueryError(ctx, err, "Hirer Profile Not Found")
			return
		}
		if registrationDate.Valid {
			date := registrationDate.Time.Format("2006-01-02")
			profile.RegistrationDate = &date
		}
		ctx.JSON(http.StatusOK, profile)
	}
}

func parseUserID(ctx *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id < 1 {
		ctx.JSON(http.StatusBadRequest, model.Fail{Error: "User ID must be a positive integer"})
		return 0, false
	}
	return id, true
}

func writeProfileQueryError(ctx *gin.Context, err error, notFoundMessage string) {
	if errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusNotFound, model.Fail{Error: notFoundMessage})
		return
	}
	ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
}
