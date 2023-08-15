package dto

import "interview_bootcamp/model"

type InterviewResultResponseDto struct {
	Id         string                 `json:"id"`
	InterviewP model.InterviewProcess `json:"interview_proses_id"`
	Result     model.Result           `json:"result_id"`
	Note       string                 `json:"note"`
}
