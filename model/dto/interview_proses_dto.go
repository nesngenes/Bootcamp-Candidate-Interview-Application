package dto

import (
	"interview_bootcamp/model"
	"time"
)

type InterviewProcessResponseDto struct {
	ID                string            `json:"id"`
	Candidate         model.Candidate   `json:"candidate"`
	Interviewer       model.Interviewer `json:"interviewer"`
	InterviewDatetime time.Time         `json:"interview_datetime"`
	MeetingLink       string            `json:"meeting_link"`
	FormInterview     string            `json:"form_interview"`
	Status            model.Status      `json:"status"`
}
