package repository

import (
	"database/sql"

	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/utils/common"
)

type StatusRepository interface {
	BaseRepository[model.Status]
	BaseRepositoryPaging[model.Status]
	GetByName(name string) (model.Status, error)
}
type statusRepository struct {
	db *sql.DB
}

func (s *statusRepository) Create(payload model.Status) error {
	_, err := s.db.Exec("INSERT INTO status (id, name) VALUES ($1, $2)", payload.StatusId, payload.Name)
	if err != nil {
		return err
	}
	return nil
}
func (s *statusRepository) GetByName(name string) (model.Status, error) {
	var status model.Status
	err := s.db.QueryRow("SELECT * FROM status WHERE name ILIKE $1", "%"+name+"%").Scan(&status.StatusId, &status.Name)
	if err != nil {
		return model.Status{}, err
	}
	return status, nil
}
func (s *statusRepository) List() ([]model.Status, error) {
	rows, err := s.db.Query("SELECT * FROM status")
	if err != nil {
		return nil, err
	}
	var statuss []model.Status
	for rows.Next() {
		var status model.Status
		err := rows.Scan(&status.StatusId, &status.Name)
		if err != nil {
			return nil, err
		}
		statuss = append(statuss, status)
	}
	return statuss, nil
}
func (s *statusRepository) Get(id string) (model.Status, error) {
	var status model.Status
	row := s.db.QueryRow("SELECT * FROM status WHERE id= $1", id)
	err := row.Scan(&status.StatusId, &status.Name)
	if err != nil {
		return model.Status{}, err
	}
	return status, nil
}
func (s *statusRepository) Update(payload model.Status) error {
	_, err := s.db.Exec("UPDATE status SET name = $2 WHERE id = $1", payload.StatusId, payload.Name)
	if err != nil {
		return err
	}
	return nil
}
func (s *statusRepository) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM status WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
func (s *statusRepository) Paging(requestPaging dto.PaginationParam) ([]model.Status, dto.Paging, error) {
	var paginationQuery dto.PaginationQuery
	paginationQuery = common.GetPaginationParams(requestPaging)
	rows, err := s.db.Query("SELECT * FROM status  LIMIT $1 OFFSET $2", paginationQuery.Take, paginationQuery.Skip)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var statuss []model.Status
	for rows.Next() {
		var status model.Status
		err := rows.Scan(&status.StatusId, &status.Name)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		statuss = append(statuss, status)
	}

	// count bootcamp
	var totalRows int
	row := s.db.QueryRow("SELECT COUNT(*) FROM status")
	err = row.Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	return statuss, common.Paginate(paginationQuery.Page, paginationQuery.Take, totalRows), nil

}
func NewStatusRepository(db *sql.DB) StatusRepository {
	return &statusRepository{db: db}
}
