package model

import "time"

type InterviewProcess struct {
	ID                string    `json:"id"`
	CandidateID       string    `json:"candidate_id"`
	InterviewerID     string    `json:"interviewer_id"`
	InterviewDatetime time.Time `json:"interview_datetime"`
	MeetingLink       string    `json:"meeting_link"`
	FormID    		  string    `json:"form_id"`
	FormInterview     string     `json:"form_interview"`
	FormLink   		  string    `json:"form_link"`
	StatusID          string    `json:"status_id"`
}
