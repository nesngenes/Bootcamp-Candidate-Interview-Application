package repository

import (
	"database/sql"
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
)

type InterviewResultRepository interface {
	Create(payload model.InterviewResult) error
	Get(id string) (dto.InterviewResultResponseDto, error)
	List(requestPaging dto.PaginationParam) ([]dto.InterviewResultResponseDto, dto.Paging, error)
}
type interviewResultRepository struct {
	db *sql.DB
}

func (i *interviewResultRepository) Create(payload model.InterviewResult) error {
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback() // Rollback the transaction if there was an error
		}
	}()

	// Insert interview Result
	_, err = tx.Exec("INSERT INTO interview_result (id, interview_id, result_id, note) VALUES ($1, $2, $3, $4)", payload.Id, payload.InterviewId, payload.ResultId, payload.Note)

	if err != nil {
		return err
	}

	err = tx.Commit() // Commit the transaction if everything is successful
	if err != nil {
		return err
	}

	return nil
}

func (i *interviewResultRepository) Get(id string) (dto.InterviewResultResponseDto, error) {
	var InterviewResultResponseDto dto.InterviewResultResponseDto

	sqlInterviewResult := `
        SELECT ir.id, ir.note ,ip.id AS InterviewP_id,ip.candidate_id AS InterviewP_candidate_id,ip.interviewer_id AS InterviewP_interviewer_id ,ip.interview_datetime AS InterviewP_interview_datetime,ip.meeting_link AS InterviewP_meeting_link,ip.form_interview AS InterviewP_form_interview,ip.status_id AS InterviewP_status_id,r.id AS result_id ,r.name AS result_name
             
        FROM interview_result AS ir
        JOIN interviews_process AS ip ON ir.interview_id = ip.id
        JOIN result AS r ON ir.result_id = r.id
        WHERE ip.id = $1`

	err := i.db.QueryRow(sqlInterviewResult, id).Scan(
		&InterviewResultResponseDto.Id, &InterviewResultResponseDto.Note,
		&InterviewResultResponseDto.InterviewP.ID,
		&InterviewResultResponseDto.InterviewP.CandidateID, &InterviewResultResponseDto.InterviewP.InterviewerID,
		&InterviewResultResponseDto.InterviewP.InterviewDatetime,
		&InterviewResultResponseDto.InterviewP.MeetingLink, &InterviewResultResponseDto.InterviewP.FormInterview,
		&InterviewResultResponseDto.InterviewP.StatusID,
		&InterviewResultResponseDto.Result.ResultId, &InterviewResultResponseDto.Result.Name,
	)
	if err != nil {
		return dto.InterviewResultResponseDto{}, err
	}

	// Add log statements to see retrieved data
	fmt.Println("Retrieved Interview Result Data:", InterviewResultResponseDto)

	return InterviewResultResponseDto, nil
}

func (i *interviewResultRepository) List(_ dto.PaginationParam) ([]dto.InterviewResultResponseDto, dto.Paging, error) {
	query := `SELECT ir.id, ir.note ,ip.id AS InterviewP_id,ip.candidate_id AS InterviewP_candidate_id,ip.interviewer_id AS InterviewP_interviewer_id ,ip.interview_datetime AS InterviewP_interview_datetime,ip.meeting_link AS InterviewP_meeting_link,ip.form_interview AS InterviewP_form_interview,ip.status_id AS InterviewP_status_id,r.id AS Result_id ,r.name AS Result_name
             
	FROM interview_result AS ir
	JOIN interviews_process AS ip ON ir.interview_id = ip.id
	JOIN result AS r ON ir.result_id = r.id`

	rows, err := i.db.Query(query)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()

	var interviewResults []dto.InterviewResultResponseDto
	for rows.Next() {
		var InterviewResult dto.InterviewResultResponseDto
		err := rows.Scan(
			&InterviewResult.Id, &InterviewResult.Note,
			&InterviewResult.InterviewP.ID,
			&InterviewResult.InterviewP.CandidateID, &InterviewResult.InterviewP.InterviewerID,
			&InterviewResult.InterviewP.InterviewDatetime,
			&InterviewResult.InterviewP.MeetingLink, &InterviewResult.InterviewP.FormInterview,
			&InterviewResult.InterviewP.StatusID,
			&InterviewResult.Result.ResultId, &InterviewResult.Result.Name,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		interviewResults = append(interviewResults, InterviewResult)
	}

	return interviewResults, dto.Paging{}, nil
}

func NewInterviewResultRepository(db *sql.DB) InterviewResultRepository {
	return &interviewResultRepository{db: db}
}
