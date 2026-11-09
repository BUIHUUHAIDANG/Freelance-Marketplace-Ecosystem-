package controller

import (
	"database/sql"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/model"
	"github.com/BUIHUUHAIDANG/freelance_market/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// LoginHandler godoc
// @Summary Sign in to a user account
// @Description Verifies the supplied email and plain-text password.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.SignInInputUserAccount true "Login credentials"
// @Success 200 {object} model.SignUpRespond
// @Failure 400 {object} model.Fail
// @Failure 401 {object} model.Fail
// @Failure 404 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /auth/login [post]
func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input model.SignInInputUserAccount
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, model.Fail{
				Error: err.Error(),
			})
			return
		}

		var userAccount model.UserAccount
		err := db.QueryRow("SELECT * FROM user_account WHERE email = ?", input.Email).Scan(&userAccount.ID, &userAccount.Email, &userAccount.PasswordHash, &userAccount.FirstName, &userAccount.LastName)
		if err == sql.ErrNoRows {
			c.JSON(404, model.Fail{
				Error: "email is not exist",
			})
			return
		}
		if err != nil {
			c.JSON(500, model.Fail{
				Error: "database query fail",
			})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(userAccount.PasswordHash), []byte(input.PasswordHash))
		if err != nil {
			c.JSON(401, model.Fail{
				Error: "Password is not correct",
			})
			return
		}
		token, err := utils.GenerateToken(int64(userAccount.ID))
		if err != nil {
			c.JSON(500, model.Fail{
				Error: err.Error(),
			})
			return
		}
		c.JSON(200, model.SignUpRespond{
			Token: "Bearer " + token,
		})
	}
}

// RegisterHandler godoc
// @Summary Register a user account
// @Description Creates an account and hashes the supplied plain-text password.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.SignUpInputUserAccount true "Account details"
// @Success 201 {object} model.UserRespond
// @Failure 400 {object} model.Fail
// @Failure 404 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /auth/register [post]
func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input model.SignUpInputUserAccount
		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, model.Fail{
				Error: err.Error(),
			})
			return
		}
		var UserAccount model.SignUpInputUserAccount
		err := db.QueryRow("SELECT * FROM user_account WHERE email = ?", input.Email).Scan(&UserAccount.Email, &UserAccount.PasswordHash, &UserAccount.FirstName, &UserAccount.LastName)
		if err != sql.ErrNoRows {
			c.JSON(404, model.Fail{
				Error: "email is exist",
			})
			return
		}
		passwordHash, err := passwordHash(input.PasswordHash)
		if err != nil {
			c.JSON(500, model.Fail{
				Error: err.Error(),
			})
			return
		}
		UserAccount.Email = input.Email
		UserAccount.PasswordHash = passwordHash
		UserAccount.FirstName = input.FirstName
		UserAccount.LastName = input.LastName
		result, err := db.Exec("INSERT INTO user_account (email,password_hash,first_name,last_name) VALUES (?,?,?,?)", UserAccount.Email, passwordHash, UserAccount.FirstName, UserAccount.LastName)
		if err != nil {
			c.JSON(500, model.Fail{
				Error: err.Error(),
			})
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			c.JSON(500, model.Fail{
				Error: err.Error(),
			})
			return
		}
		c.JSON(201, model.UserRespond{
			ID:        id,
			Email:     UserAccount.Email,
			FirstName: UserAccount.FirstName,
			LastName:  UserAccount.LastName,
		})
	}
}
