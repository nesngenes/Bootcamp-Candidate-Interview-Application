package model

type HRRecruitment struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	UserID   string `json:"user_id"`
	User     Users  `json:"user"`
}
