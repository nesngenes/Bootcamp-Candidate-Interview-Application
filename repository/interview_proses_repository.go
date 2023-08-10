package repository

import (
	"database/sql"
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
	// insert interviewProces
	_, err = tx.Exec("INSERT INTO interviews_process (id, candidate_id, interviewer_id, interviw_datetime, meeting_link ,from_interview, status_id) VALUES ($1, $2, $3, $4, $5, $6,$7)", payload.ID, payload.CandidateID, payload.InterviewerID, payload.InterviewDatetime, payload.MeetingLink, payload.FormInterview, payload.StatusID)

	if err != nil {
		return err
	}
	return nil
}
func (i *interviewProcessRepository) Get(id string) (dto.InterviewProcessResponseDto, error) {
	var interviewPrResponseDto dto.InterviewProcessResponseDto
	sqlInterviewProcess := `SELECT ip.id, ip.interview_datetime, ip.meeting_link, ip.from_interview, c.id AS candidate_id, c.full_name AS candidate_full_name, c.phone_number AS candidate_phone, c.email AS candidate_email, c.date_of_birth AS candidate_date_of_birth, c.address AS candidate_address, c.cv_link AS candidate_cv_link, c.bootcamp_id AS candidate_bootcamp_id, c.instansi_pendidikan AS candidate_instansi_pendidikan, c.hacker_rank AS candidate_hacker_rank, i.id AS interviewer_id, i.full_name AS interviewer_full_name, i.user_id AS interviewer_user_id, s.id AS status_id, s.name AS status_name
	FROM
    interviews_process AS ip
    JOIN
    candidate AS c ON ip.candidate_id = c.id
    JOIN
    interviewer AS i ON ip.interviewer_id = i.id
    JOIN
    status AS s ON ip.status_id = s.id
    WHERE
    ip.id = $1`

	err := i.db.QueryRow(sqlInterviewProcess, id).Scan(&interviewPrResponseDto.ID, &interviewPrResponseDto.Candidate.CandidateID, &interviewPrResponseDto.Candidate.FullName, &interviewPrResponseDto.Candidate.Phone, &interviewPrResponseDto.Candidate.Email, &interviewPrResponseDto.Candidate.DateOfBirth, &interviewPrResponseDto.Candidate.Address, &interviewPrResponseDto.Candidate.CvLink, &interviewPrResponseDto.Candidate.Bootcamp.BootcampId, &interviewPrResponseDto.Candidate.InstansiPendidikan, &interviewPrResponseDto.Candidate.HackerRank, &interviewPrResponseDto.Interviewer.InterviewerID, &interviewPrResponseDto.Interviewer.FullName, &interviewPrResponseDto.Interviewer.UserID, &interviewPrResponseDto.InterviewDatetime, &interviewPrResponseDto.MeetingLink, &interviewPrResponseDto.FormInterview, &interviewPrResponseDto.Status.StatusId, &interviewPrResponseDto.Status.Name)
	if err != nil {
		return dto.InterviewProcessResponseDto{}, err
	}

	return interviewPrResponseDto, nil
}
func (i *interviewProcessRepository) List(requestPaging dto.PaginationParam) ([]dto.InterviewProcessResponseDto, dto.Paging, error) {
	return nil, dto.Paging{}, nil
}

func NewInterviewProcessRepository(db *sql.DB) InterviewProcessRepository {
	return &interviewProcessRepository{db: db}
}
