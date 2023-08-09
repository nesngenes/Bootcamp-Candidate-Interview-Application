package model

type Interviewer struct {
	InterviewerID  string `json:"interviewer_id"`
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Specialization string `json:"specialization"`
	RoleName       string `json:"role_name"`
}
