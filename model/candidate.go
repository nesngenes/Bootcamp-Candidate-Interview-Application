package model

type Candidate struct {
	CandidateID        string   `json:"candidate_id"`
	FullName           string   `json:"full_name"`
	Phone              string   `json:"phone"`
	Email              string   `json:"email"`
	DateOfBirth        string   `json:"date_of_birth"`
	Address            string   `json:"address"`
	CvLink             string   `json:"cv_link"`
	Bootcamp           Bootcamp `json:"bootcamp"`
	InstansiPendidikan string   `json:"instansi_pendidikan"`
	HackerRank         int      `json:"hackerrank_score"`
}
