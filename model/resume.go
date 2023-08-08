package model

type Resume struct {
	ResumeID    string `json:"resume_id"`
	CandidateID string `json:"candidate_id"`
	CvURL       string `json:"cv_url"`
	CvFile      []byte `json:"-"` // This field will not be serialized to JSON
	CvFileName  string `json:"-"`
}
