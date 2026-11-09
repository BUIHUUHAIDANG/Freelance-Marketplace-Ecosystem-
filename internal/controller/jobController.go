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
	"github.com/go-sql-driver/mysql"
)

func createInClausePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("?,", n-1) + "?"
}

// GetAllJobs godoc
// @Summary Get all jobs
// @Description Returns a paginated list of jobs with their required skills. Invalid page or limit values use defaults.
// @Tags jobs
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Jobs per page (maximum 10)" default(10)
// @Success 200 {object} model.JobListRespond
// @Failure 500 {object} model.Fail
// @Router /job [get]
func GetAllJobs(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil || limit > 10 || limit < 1 {
			limit = 10
		}
		offset := (page - 1) * limit
		query := `SELECT j.id,u.first_name,e.duration_text,c.complexity_type,j.description, s.skill_name,p.type_name,j.payment_amount,j.title FROM job j 
				INNER JOIN hire_manager h ON j.hire_manager_id  = h.id
				INNER JOIN user_account u  ON u.id = h.user_account_id
				INNER JOIN expected_duration e  ON j.duration_id = e.id
				INNER JOIN complexity c ON j.complexity_id  = c.id
				INNER JOIN job_require_skill j2 ON j2.job_id  = j.id
				INNER JOIN skill s ON j2.skill_id  = s.id
				INNER JOIN payment_type p ON p.id  =j.payment_type_id 
				WHERE j.main_skill_id  = s.id
				LIMIT ?
				OFFSET ?
		`
		rows, err := db.QueryContext(ctx.Request.Context(), query, limit, offset)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{
				Error: err.Error(),
			})
			return
		}
		listJobs := make([]model.Job, 0)
		listJobsId := make([]int64, 0)
		jobMap := make(map[int64]*model.Job)

		for rows.Next() {
			var j model.Job
			j.RequireSkill = make([]model.OtherSkill, 0)
			err := rows.Scan(
				&j.ID,
				&j.HireManagerName,
				&j.Duration,
				&j.Complexity,
				&j.Description,
				&j.MainSkill,
				&j.PaymentType,
				&j.PaymentAmount,
				&j.Title,
			)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, model.Fail{
					Error: err.Error(),
				})
				return
			}
			listJobs = append(listJobs, j)
			listJobsId = append(listJobsId, j.ID)
		}
		if err := rows.Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{
				Error: err.Error(),
			})
			return
		}
		if len(listJobs) > 0 {
			querySkills := fmt.Sprintf(`
				SELECT jrs.job_id, s.skill_name 
				FROM job_require_skill jrs
				INNER JOIN skill s ON jrs.skill_id = s.id
				WHERE jrs.job_id IN (%s)
			`, createInClausePlaceholders(len(listJobsId)))

			args := make([]interface{}, len(listJobsId))
			for i, id := range listJobsId {
				args[i] = id
			}

			// Map lại địa chỉ con trỏ để gán skill cho từng job nhanh chóng
			for i := range listJobs {
				jobMap[listJobs[i].ID] = &listJobs[i]
			}

			skillRows, err := db.QueryContext(ctx.Request.Context(), querySkills, args...)
			if err == nil {
				defer skillRows.Close()
				for skillRows.Next() {
					var jobID int64
					var skillName string
					if err := skillRows.Scan(&jobID, &skillName); err == nil {
						if job, exists := jobMap[jobID]; exists {
							job.RequireSkill = append(job.RequireSkill, model.OtherSkill{
								SkillName: skillName,
							})
						}
					}
				}
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"page":  page,
			"limit": limit,
			"data":  listJobs,
		})

	}
}

// GetJobById godoc
// @Summary Get a job by ID
// @Description Returns a job and its required skills.
// @Tags jobs
// @Produce json
// @Param jobid path int true "Job ID"
// @Success 200 {object} model.JobRespond
// @Failure 400 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /job/{jobid} [get]
func GetJobById(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jobId := ctx.Param("jobid")
		var j model.Job
		j.RequireSkill = make([]model.OtherSkill, 0)
		query := `
		SELECT j.id,u.first_name,e.duration_text,c.complexity_type,j.description, s.skill_name,p.type_name,j.payment_amount,j.title FROM job j 
		INNER JOIN hire_manager h ON j.hire_manager_id  = h.id
		INNER JOIN user_account u  ON u.id = h.user_account_id
		INNER JOIN expected_duration e  ON j.duration_id = e.id
		INNER JOIN complexity c ON j.complexity_id  = c.id
		INNER JOIN job_require_skill j2 ON j2.job_id  = j.id
		INNER JOIN skill s ON j2.skill_id  = s.id
		INNER JOIN payment_type p ON p.id  =j.payment_type_id 
		WHERE j.main_skill_id  = s.id AND j.id = ?
		`
		err := db.QueryRow(query, jobId).Scan(
			&j.ID,
			&j.HireManagerName,
			&j.Duration,
			&j.Complexity,
			&j.Description,
			&j.MainSkill,
			&j.PaymentType,
			&j.PaymentAmount,
			&j.Title,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(400, model.Fail{
					Error: "Job not exist",
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, model.Fail{
				Error: err.Error(),
			})
			return
		}
		query2 := `
		SELECT  s.skill_name 
				FROM job_require_skill jrs
				INNER JOIN skill s ON jrs.skill_id = s.id
				WHERE jrs.job_id = ?
		`
		rowsSkill, err := db.QueryContext(ctx.Request.Context(), query2, jobId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{
				Error: err.Error(),
			})
			return
		}
		defer rowsSkill.Close()
		for rowsSkill.Next() {
			var skillNameScan string
			err := rowsSkill.Scan(&skillNameScan)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, model.Fail{
					Error: err.Error(),
				})
				return
			}
			j.RequireSkill = append(j.RequireSkill, model.OtherSkill{
				SkillName: skillNameScan,
			})
		}
		ctx.JSON(200, gin.H{
			"data": j,
		})
	}
}

// AddJob godoc
// @Summary Create a job
// @Description Hire Managers can create jobs for their own profile. Admin user IDs configured in ADMIN_USER_IDS can create jobs for any manager.
// @Tags jobs
// @Accept json
// @Produce json
// @Param request body model.AddJobRequest true "Job details"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} model.Fail
// @Failure 403 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /job [post]
func AddJob(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.AddJobRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Fail{Error: err.Error()})
			return
		}
		if err := validateAddJobRequest(input); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Fail{Error: err.Error()})
			return
		}
		if !isAdminRequest(ctx) {
			hireManagerID, ok := getAuthenticatedHireManagerID(ctx, db)
			if !ok {
				return
			}
			if input.HireManagerID != hireManagerID {
				ctx.JSON(http.StatusForbidden, model.Fail{Error: "You can only create jobs for your own hire manager profile"})
				return
			}
		}

		tx, err := db.BeginTx(ctx.Request.Context(), nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}
		defer tx.Rollback()

		result, err := tx.ExecContext(ctx.Request.Context(), `
			INSERT INTO job
				(hire_manager_id, duration_id, complexity_id, description,
				 main_skill_id, payment_type_id, payment_amount, title)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			input.HireManagerID, input.DurationID, input.ComplexityID, input.Description,
			input.MainSkillID, input.PaymentTypeID, input.PaymentAmount, strings.TrimSpace(input.Title))
		if err != nil {
			writeJobMutationError(ctx, err)
			return
		}
		jobID, err := result.LastInsertId()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}
		if err := replaceJobSkills(ctx, tx, jobID, input); err != nil {
			writeJobMutationError(ctx, err)
			return
		}
		if err := tx.Commit(); err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Job created successfully",
			"job_id":  jobID,
		})
	}
}

// UpdateJobById godoc
// @Summary Update a job
// @Description Hire Managers can update only their own jobs. Admin user IDs configured in ADMIN_USER_IDS can update any job.
// @Tags jobs
// @Accept json
// @Produce json
// @Param jobid path int true "Job ID"
// @Param request body model.AddJobRequest true "Updated job details"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} model.Fail
// @Failure 403 {object} model.Fail
// @Failure 404 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /job/{jobid} [put]
func UpdateJobById(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jobID, ok := parseJobID(ctx)
		if !ok {
			return
		}
		var input model.AddJobRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Fail{Error: err.Error()})
			return
		}
		if err := validateAddJobRequest(input); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Fail{Error: err.Error()})
			return
		}
		ownerHireManagerID, ok := authorizeJobManager(ctx, db, jobID)
		if !ok {
			return
		}
		if !isAdminRequest(ctx) && input.HireManagerID != ownerHireManagerID {
			ctx.JSON(http.StatusForbidden, model.Fail{Error: "You cannot transfer a job to another hire manager"})
			return
		}

		tx, err := db.BeginTx(ctx.Request.Context(), nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}
		defer tx.Rollback()

		result, err := tx.ExecContext(ctx.Request.Context(), `
			UPDATE job
			SET hire_manager_id = ?, duration_id = ?, complexity_id = ?, description = ?,
			    main_skill_id = ?, payment_type_id = ?, payment_amount = ?, title = ?
			WHERE id = ?`,
			input.HireManagerID, input.DurationID, input.ComplexityID, input.Description,
			input.MainSkillID, input.PaymentTypeID, input.PaymentAmount, strings.TrimSpace(input.Title), jobID)
		if err != nil {
			writeJobMutationError(ctx, err)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}
		if rowsAffected == 0 {
			var exists bool
			if err := tx.QueryRowContext(ctx.Request.Context(),
				"SELECT EXISTS(SELECT 1 FROM job WHERE id = ?)", jobID).Scan(&exists); err != nil {
				ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
				return
			}
			if !exists {
				ctx.JSON(http.StatusNotFound, model.Fail{Error: "Job not found"})
				return
			}
		}

		if err := replaceJobSkills(ctx, tx, jobID, input); err != nil {
			writeJobMutationError(ctx, err)
			return
		}
		if err := tx.Commit(); err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Job updated successfully",
			"job_id":  jobID,
		})
	}
}

// DeleteJobById godoc
// @Summary Delete a job
// @Description Hire Managers can delete only their own jobs. Admin user IDs configured in ADMIN_USER_IDS can delete any job.
// @Tags jobs
// @Produce json
// @Param jobid path int true "Job ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} model.Fail
// @Failure 403 {object} model.Fail
// @Failure 404 {object} model.Fail
// @Failure 409 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /job/{jobid} [delete]
func DeleteJobById(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jobID, ok := parseJobID(ctx)
		if !ok {
			return
		}
		if _, ok := authorizeJobManager(ctx, db, jobID); !ok {
			return
		}
		result, err := db.ExecContext(ctx.Request.Context(), "DELETE FROM job WHERE id = ?", jobID)
		if err != nil {
			writeJobMutationError(ctx, err)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}
		if rowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, model.Fail{Error: "Job not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
	}
}

// ApplyJob godoc
// @Summary Apply for a job
// @Description Creates a proposal for the authenticated freelancer. Admin user IDs configured in ADMIN_USER_IDS may also apply if they have a freelancer profile. The proposal starts with Pending status.
// @Tags jobs
// @Accept json
// @Produce json
// @Param jobid path int true "Job ID"
// @Param request body model.ApplyJobRequest true "Proposal details"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} model.Fail
// @Failure 403 {object} model.Fail
// @Failure 404 {object} model.Fail
// @Failure 409 {object} model.Fail
// @Failure 500 {object} model.Fail
// @Router /job/{jobid}/apply [post]
func ApplyJob(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jobID, ok := parseJobID(ctx)
		if !ok {
			return
		}
		var input model.ApplyJobRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Fail{Error: err.Error()})
			return
		}
		if input.PaymentTypeID < 1 || input.PaymentAmount <= 0 {
			ctx.JSON(http.StatusBadRequest, model.Fail{Error: "payment_type_id and a positive payment_amount are required"})
			return
		}
		userID, ok := authenticatedUserID(ctx)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, model.Fail{Error: "Unauthorized"})
			return
		}

		var freelancerID int64
		err := db.QueryRowContext(ctx.Request.Context(),
			"SELECT id FROM freelancer WHERE user_account_id = ?", userID).Scan(&freelancerID)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusForbidden, model.Fail{Error: "A freelancer profile is required to apply for a job"})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}

		tx, err := db.BeginTx(ctx.Request.Context(), nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}
		defer tx.Rollback()

		var existingJobID int64
		err = tx.QueryRowContext(ctx.Request.Context(),
			"SELECT id FROM job WHERE id = ? FOR UPDATE", jobID).Scan(&existingJobID)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, model.Fail{Error: "Job not found"})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}

		var alreadyApplied bool
		if err := tx.QueryRowContext(ctx.Request.Context(),
			"SELECT EXISTS(SELECT 1 FROM proposal WHERE job_id = ? AND freelancer_id = ?)",
			jobID, freelancerID).Scan(&alreadyApplied); err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}
		if alreadyApplied {
			ctx.JSON(http.StatusConflict, model.Fail{Error: "You have already applied for this job"})
			return
		}

		var pendingStatusID int64
		err = tx.QueryRowContext(ctx.Request.Context(),
			"SELECT id FROM status_catalog WHERE status_name = ?", "Pending").Scan(&pendingStatusID)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: "Pending status is missing from status_catalog"})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}

		result, err := tx.ExecContext(ctx.Request.Context(), `
			INSERT INTO proposal
				(job_id, freelancer_id, payment_type_id, payment_amount,
				 current_proposal_status_id, freelancer_comment)
			VALUES (?, ?, ?, ?, ?, ?)`,
			jobID, freelancerID, input.PaymentTypeID, input.PaymentAmount,
			pendingStatusID, input.FreelancerComment)
		if err != nil {
			writeJobMutationError(ctx, err)
			return
		}
		proposalID, err := result.LastInsertId()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}
		if err := tx.Commit(); err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message":     "Application submitted successfully",
			"proposal_id": proposalID,
			"job_id":      jobID,
			"status":      "Pending",
		})
	}
}

func parseJobID(ctx *gin.Context) (int64, bool) {
	jobID, err := strconv.ParseInt(ctx.Param("jobid"), 10, 64)
	if err != nil || jobID < 1 {
		ctx.JSON(http.StatusBadRequest, model.Fail{Error: "Job ID must be a positive integer"})
		return 0, false
	}
	return jobID, true
}

func validateAddJobRequest(input model.AddJobRequest) error {
	if input.HireManagerID < 1 || input.DurationID < 1 || input.ComplexityID < 1 ||
		input.MainSkillID < 1 || input.PaymentTypeID < 1 {
		return errors.New("all job reference IDs must be positive integers")
	}
	if strings.TrimSpace(input.Title) == "" {
		return errors.New("title is required")
	}
	if input.PaymentAmount != nil && *input.PaymentAmount < 0 {
		return errors.New("payment_amount cannot be negative")
	}
	for _, skillID := range input.RequiredSkillIDs {
		if skillID < 1 {
			return errors.New("required_skill_ids must contain positive integers")
		}
	}
	return nil
}

func replaceJobSkills(ctx *gin.Context, tx *sql.Tx, jobID int64, input model.AddJobRequest) error {
	if _, err := tx.ExecContext(ctx.Request.Context(), "DELETE FROM job_require_skill WHERE job_id = ?", jobID); err != nil {
		return err
	}

	skillIDs := make([]int64, 0, len(input.RequiredSkillIDs)+1)
	seen := make(map[int64]struct{}, len(input.RequiredSkillIDs)+1)
	skillIDs = append(skillIDs, input.MainSkillID)
	seen[input.MainSkillID] = struct{}{}
	for _, skillID := range input.RequiredSkillIDs {
		if _, exists := seen[skillID]; exists {
			continue
		}
		seen[skillID] = struct{}{}
		skillIDs = append(skillIDs, skillID)
	}

	for _, skillID := range skillIDs {
		if _, err := tx.ExecContext(ctx.Request.Context(),
			"INSERT INTO job_require_skill (job_id, skill_id) VALUES (?, ?)", jobID, skillID); err != nil {
			return err
		}
	}
	return nil
}

func writeJobMutationError(ctx *gin.Context, err error) {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1452:
			ctx.JSON(http.StatusBadRequest, model.Fail{Error: "A referenced job record does not exist"})
			return
		case 1451:
			ctx.JSON(http.StatusConflict, model.Fail{Error: "Job is referenced by other records and cannot be deleted"})
			return
		}
	}
	ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
}

func authenticatedUserID(ctx *gin.Context) (int64, bool) {
	value, exists := ctx.Get("user_id")
	userID, ok := value.(int64)
	return userID, exists && ok && userID > 0
}

func isAdminRequest(ctx *gin.Context) bool {
	value, _ := ctx.Get("is_admin")
	isAdmin, _ := value.(bool)
	return isAdmin
}

func getAuthenticatedHireManagerID(ctx *gin.Context, db *sql.DB) (int64, bool) {
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.Fail{Error: "Unauthorized"})
		return 0, false
	}
	var hireManagerID int64
	err := db.QueryRowContext(ctx.Request.Context(),
		"SELECT id FROM hire_manager WHERE user_account_id = ?", userID).Scan(&hireManagerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
		return 0, false
	}
	return hireManagerID, true
}

func authorizeJobManager(ctx *gin.Context, db *sql.DB, jobID int64) (int64, bool) {
	var hireManagerID int64
	err := db.QueryRowContext(ctx.Request.Context(),
		"SELECT hire_manager_id FROM job WHERE id = ?", jobID).Scan(&hireManagerID)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusNotFound, model.Fail{Error: "Job not found"})
		return 0, false
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
		return 0, false
	}
	if isAdminRequest(ctx) {
		return hireManagerID, true
	}

	userID, ok := authenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.Fail{Error: "Unauthorized"})
		return 0, false
	}
	var ownsJob bool
	err = db.QueryRowContext(ctx.Request.Context(), `
		SELECT EXISTS(
			SELECT 1 FROM hire_manager
			WHERE id = ? AND user_account_id = ?
		)`, hireManagerID, userID).Scan(&ownsJob)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Fail{Error: err.Error()})
		return 0, false
	}
	if !ownsJob {
		ctx.JSON(http.StatusForbidden, model.Fail{Error: "You can only manage your own jobs"})
		return 0, false
	}
	return hireManagerID, true
}
