package model


type Config struct {
	Port string
}
// UserAccount represents a user account in the system.
type SignUpInputUserAccount struct {
	Email    string `json:"email" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
}

type SignInInputUserAccount struct {
	Email    string `json:"email" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
}

type UserRespond struct {
	ID        int64    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}


type UserAccount struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Fail struct{
	Error string `json:"error"`
}


type UserAccountProfileRespond struct{
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	HasFreelancerProfile bool `json:"has_freelancer_profile"`
	HasHirerProfile bool `json:"has_hirer_profile"`
}