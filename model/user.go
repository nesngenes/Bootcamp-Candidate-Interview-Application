package model

type Users struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	UserName string   `json:"user_name"`
	Password string   `json:"password"`
	UserRole UserRole `json:"user_role"`
}
