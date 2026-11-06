package controller

import (
	"database/sql"

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
// @Failure 404 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /user/{id} [get]
func GetUserProfile(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userAccountProfile model.UserAccountProfileRespond
		userId := ctx.Param("id")
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
		err := db.QueryRowContext(ctx.Request.Context(), query, userId).
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
