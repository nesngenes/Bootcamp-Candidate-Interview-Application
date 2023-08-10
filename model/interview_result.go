package model

type InterviewResult struct {
	Id          string `json:"id"`
	InterviewId string `json:"interview_id"`
	ResultId    string `json:"result_id"`
	Note        string `json:"note"`
}
