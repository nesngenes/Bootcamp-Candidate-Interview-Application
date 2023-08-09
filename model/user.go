package model

type Users struct {
	Id       string    `json:"id"`
	Email    string    `json:"email"`
	UserName string    `json:"user_name"`
	Password string    `json:"password"`
	UserRole UserRoles `json:"user_role"`
}
