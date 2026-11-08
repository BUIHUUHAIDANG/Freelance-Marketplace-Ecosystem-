package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/model"
	"github.com/gin-gonic/gin"
)


func createInClausePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("?,", n-1) + "?"
}

func GetAllJobs(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page,err := strconv.Atoi(ctx.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		limit,err := strconv.Atoi(ctx.Query("limit"))
		if err != nil || limit > 10 || limit < 1 {
			limit = 10
		}
		offset := (page -1) *limit
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
		rows, err := db.QueryContext(ctx.Request.Context(),query,limit,offset)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,model.Fail{
				Error: err.Error(),
			})
			return 
		}
		listJobs := make([]model.Job,0)
		listJobsId := make([]int64,0)
		jobMap := make(map[int64]*model.Job)

		for rows.Next(){
			var j model.Job
			j.RequireSkill = make([]model.OtherSkill,0)
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
				ctx.JSON(http.StatusInternalServerError,model.Fail{
					Error: err.Error(),
				})
				return 
			}
			listJobs = append(listJobs, j)
			listJobsId =append(listJobsId, j.ID)
		}
		if err := rows.Err(); err !=nil{
			ctx.JSON(http.StatusInternalServerError,model.Fail{
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
		ctx.JSON(http.StatusOK,gin.H{
			"page":  page,
			"limit": limit,
			"data":  listJobs,
		})
		
	}
}

func GetJobById(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jobId := ctx.Param("jobid")
		var j model.Job
		j.RequireSkill = make([]model.OtherSkill,0)
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
		err := db.QueryRow(query,jobId).Scan(
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
			if err == sql.ErrNoRows{
				ctx.JSON(400,model.Fail{
					Error: "Job not exist",
				})
				return 
			}
			ctx.JSON(http.StatusInternalServerError,model.Fail{
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
		rowsSkill,err := db.QueryContext(ctx.Request.Context(),query2,jobId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,model.Fail{
				Error: err.Error(),
			})
			return 
		}		
		defer rowsSkill.Close()
		for rowsSkill.Next(){
			var skillNameScan string
			err := rowsSkill.Scan(&skillNameScan)
			if err != nil{
				ctx.JSON(http.StatusInternalServerError,model.Fail{
					Error: err.Error(),
				})
				return 
			}
			j.RequireSkill = append(j.RequireSkill, model.OtherSkill{
				SkillName: skillNameScan,
			})
		}
		ctx.JSON(200,gin.H{
			"data":  j,
		})	
	}
}


func AddJob(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}


func UpdateJobById(db *sql.DB) gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func DeleteJobById(db *sql.DB) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
	}
}