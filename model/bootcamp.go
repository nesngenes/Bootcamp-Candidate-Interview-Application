package model

import "time"

type Bootcamp struct {
	BootcampId string    `json:"id"`
	Name       string    `json:"name"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Location   string    `json:"location"`
	Id        string    `json:"id"`
}
