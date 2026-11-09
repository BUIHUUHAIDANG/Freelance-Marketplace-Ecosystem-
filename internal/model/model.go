package model

type Config struct {
	Port string
}

// UserAccount represents a user account in the system.
type SignUpInputUserAccount struct {
	Email        string `json:"email" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
}

type SignInInputUserAccount struct {
	Email        string `json:"email" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
}

type SignUpRespond struct {
	Token string `json:"token"`
}

type UserRespond struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserAccount struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type Fail struct {
	Error string `json:"error"`
}

type UserAccountProfileRespond struct {
	Email                string `json:"email"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	HasFreelancerProfile bool   `json:"has_freelancer_profile"`
	HasHirerProfile      bool   `json:"has_hirer_profile"`
}

type UpdateUserProfileRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

type FreelancerProfileRespond struct {
	ID               int64   `json:"id"`
	UserAccountID    int64   `json:"user_account_id"`
	RegistrationDate *string `json:"registration_date"`
	Location         *string `json:"location"`
	Overview         *string `json:"overview"`
}

type HirerProfileRespond struct {
	ID               int64   `json:"id"`
	UserAccountID    int64   `json:"user_account_id"`
	RegistrationDate *string `json:"registration_date"`
	Location         *string `json:"location"`
	CompanyID        int64   `json:"company_id"`
	CompanyName      string  `json:"company_name"`
	CompanyLocation  *string `json:"company_location"`
}

// job
type Job struct {
	ID              int64        `json:"id"`
	HireManagerName string       `json:"hire_manager_name"`
	Duration        string       `json:"duration"`
	Complexity      string       `json:"complexity"`
	Description     *string      `json:"description"`
	MainSkill       string       `json:"main_skill"`
	PaymentType     string       `json:"payment_type"`
	PaymentAmount   *float64     `json:"payment_amount"`
	Title           string       `json:"title"`
	RequireSkill    []OtherSkill `json:"required_skill"`
}

type JobListRespond struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Data  []Job `json:"data"`
}

type JobRespond struct {
	Data Job `json:"data"`
}

type OtherSkill struct {
	SkillName string `json:"other_skill"`
}

type AddJobRequest struct {
	HireManagerID    int64    `json:"hire_manager_id" binding:"required"`
	DurationID       int64    `json:"duration_id" binding:"required"`
	ComplexityID     int64    `json:"complexity_id" binding:"required"`
	Description      *string  `json:"description"`
	MainSkillID      int64    `json:"main_skill_id" binding:"required"`
	PaymentTypeID    int64    `json:"payment_type_id" binding:"required"`
	PaymentAmount    *float64 `json:"payment_amount"`
	Title            string   `json:"title" binding:"required"`
	RequiredSkillIDs []int64  `json:"required_skill_ids"`
}

type ApplyJobRequest struct {
	PaymentTypeID     int64   `json:"payment_type_id" binding:"required"`
	PaymentAmount     float64 `json:"payment_amount" binding:"required"`
	FreelancerComment *string `json:"freelancer_comment"`
}
