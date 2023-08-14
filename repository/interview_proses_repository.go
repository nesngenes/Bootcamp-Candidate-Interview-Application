package repository

import (
	"database/sql"
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
)

type InterviewProcessRepository interface {
	Create(payload model.InterviewProcess) error
	Get(id string) (dto.InterviewProcessResponseDto, error)
	List(requestPaging dto.PaginationParam) ([]dto.InterviewProcessResponseDto, dto.Paging, error)
}
type interviewProcessRepository struct {
	db *sql.DB
}

func (i *interviewProcessRepository) Create(payload model.InterviewProcess) error {
    tx, err := i.db.Begin()
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback() // Rollback the transaction if there was an error
        }
    }()

    // Insert interview process
    _, err = tx.Exec("INSERT INTO interviews_process (id, candidate_id, interviewer_id, interview_datetime, meeting_link, form_interview, status_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", payload.ID, payload.CandidateID, payload.InterviewerID, payload.InterviewDatetime, payload.MeetingLink, payload.FormLink, payload.StatusID)

    if err != nil {
        return err
    }

    err = tx.Commit() // Commit the transaction if everything is successful
    if err != nil {
        return err
    }

    return nil
}


func (i *interviewProcessRepository) Get(id string) (dto.InterviewProcessResponseDto, error) {
	var interviewPrResponseDto dto.InterviewProcessResponseDto

	sqlInterviewProcess := `
        SELECT ip.id, ip.interview_datetime, ip.meeting_link, ip.form_interview,
               c.id AS candidate_id, c.full_name AS candidate_full_name,
               c.email AS candidate_email, c.date_of_birth AS candidate_date_of_birth,
               c.address AS candidate_address, c.cv_link AS candidate_cv_link,
               c.bootcamp_id AS candidate_bootcamp_id, c.instansi_pendidikan AS candidate_instansi_pendidikan,
               c.hackerrank_score AS candidate_hacker_rank,
               i.id AS interviewer_id, i.full_name AS interviewer_full_name, i.user_id AS interviewer_user_id,
               s.id AS status_id, s.name AS status_name
        FROM interviews_process AS ip
        JOIN candidate AS c ON ip.candidate_id = c.id
        JOIN interviewer AS i ON ip.interviewer_id = i.id
        JOIN status AS s ON ip.status_id = s.id
        WHERE ip.id = $1`

	err := i.db.QueryRow(sqlInterviewProcess, id).Scan(
		&interviewPrResponseDto.ID, &interviewPrResponseDto.InterviewDatetime, &interviewPrResponseDto.MeetingLink,
		&interviewPrResponseDto.FormLink,
		&interviewPrResponseDto.Candidate.CandidateID, &interviewPrResponseDto.Candidate.FullName,
		&interviewPrResponseDto.Candidate.Email,
		&interviewPrResponseDto.Candidate.DateOfBirth, &interviewPrResponseDto.Candidate.Address,
		&interviewPrResponseDto.Candidate.CvLink, &interviewPrResponseDto.Candidate.Bootcamp.BootcampId,
		&interviewPrResponseDto.Candidate.InstansiPendidikan, &interviewPrResponseDto.Candidate.HackerRank,
		&interviewPrResponseDto.Interviewer.InterviewerID, &interviewPrResponseDto.Interviewer.FullName,
		&interviewPrResponseDto.Interviewer.UserID,
		&interviewPrResponseDto.Status.StatusId, &interviewPrResponseDto.Status.Name,
	)
	if err != nil {
		return dto.InterviewProcessResponseDto{}, err
	}

	// Add log statements to see retrieved data
	fmt.Println("Retrieved Interview Process Data:", interviewPrResponseDto)

	return interviewPrResponseDto, nil
}

func (i *interviewProcessRepository) List(_ dto.PaginationParam) ([]dto.InterviewProcessResponseDto, dto.Paging, error) {
	query := `
        SELECT ip.id, ip.interview_datetime, ip.meeting_link, ip.form_interview,
               c.id AS candidate_id, c.full_name AS candidate_full_name,
               c.email AS candidate_email, c.date_of_birth AS candidate_date_of_birth,
               c.address AS candidate_address, c.cv_link AS candidate_cv_link,
               c.bootcamp_id AS candidate_bootcamp_id, c.instansi_pendidikan AS candidate_instansi_pendidikan,
               c.hackerrank_score AS candidate_hacker_rank,
               i.id AS interviewer_id, i.full_name AS interviewer_full_name, i.user_id AS interviewer_user_id,
               s.id AS status_id, s.name AS status_name
        FROM interviews_process AS ip
        JOIN candidate AS c ON ip.candidate_id = c.id
        JOIN interviewer AS i ON ip.interviewer_id = i.id
        JOIN status AS s ON ip.status_id = s.id
    `

	rows, err := i.db.Query(query)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()

	var interviewProcesses []dto.InterviewProcessResponseDto
	for rows.Next() {
		var interviewProcess dto.InterviewProcessResponseDto
		err := rows.Scan(
			&interviewProcess.ID, &interviewProcess.InterviewDatetime, &interviewProcess.MeetingLink,
			&interviewProcess.FormLink,
			&interviewProcess.Candidate.CandidateID, &interviewProcess.Candidate.FullName,
			&interviewProcess.Candidate.Email,
			&interviewProcess.Candidate.DateOfBirth, &interviewProcess.Candidate.Address,
			&interviewProcess.Candidate.CvLink, &interviewProcess.Candidate.Bootcamp.BootcampId,
			&interviewProcess.Candidate.InstansiPendidikan, &interviewProcess.Candidate.HackerRank,
			&interviewProcess.Interviewer.InterviewerID, &interviewProcess.Interviewer.FullName,
			&interviewProcess.Interviewer.UserID,
			&interviewProcess.Status.StatusId, &interviewProcess.Status.Name,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		interviewProcesses = append(interviewProcesses, interviewProcess)
	}

	return interviewProcesses, dto.Paging{}, nil
}

func NewInterviewProcessRepository(db *sql.DB) InterviewProcessRepository {
	return &interviewProcessRepository{db: db}
}
