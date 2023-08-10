package repository

import (
	"database/sql"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/utils/common"
)

type ResultRepository interface {
	BaseRepository[model.Result]
	BaseRepositoryPaging[model.Result]
	GetByName(name string) (model.Result, error)
}
type resultRepository struct {
	db *sql.DB
}

func (r *resultRepository) Create(payload model.Result) error {
	_, err := r.db.Exec("INSERT INTO result (id, name) VALUES ($1, $2)", payload.ResultId, payload.Name)
	if err != nil {
		return err
	}
	return nil
}
func (r *resultRepository) GetByName(name string) (model.Result, error) {
	var result model.Result
	err := r.db.QueryRow("SELECT * FROM result WHERE name ILIKE $1", "%"+name+"%").Scan(&result.ResultId, &result.Name)
	if err != nil {
		return model.Result{}, err
	}
	return result, nil
}
func (r *resultRepository) List() ([]model.Result, error) {
	rows, err := r.db.Query("SELECT * FROM result")
	if err != nil {
		return nil, err
	}
	var results []model.Result
	for rows.Next() {
		var result model.Result
		err := rows.Scan(&result.ResultId, &result.Name)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
func (r *resultRepository) Get(id string) (model.Result, error) {
	var result model.Result
	row := r.db.QueryRow("SELECT * FROM result WHERE id= $1", id)
	err := row.Scan(&result.ResultId, &result.Name)
	if err != nil {
		return model.Result{}, err
	}
	return result, nil
}
func (r *resultRepository) Update(payload model.Result) error {
	_, err := r.db.Exec("UPDATE result SET name = $2 WHERE id = $1", payload.ResultId, payload.Name)
	if err != nil {
		return err
	}
	return nil
}
func (r *resultRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM result WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
func (r *resultRepository) Paging(requestPaging dto.PaginationParam) ([]model.Result, dto.Paging, error) {
	var paginationQuery dto.PaginationQuery
	paginationQuery = common.GetPaginationParams(requestPaging)
	rows, err := r.db.Query("SELECT * FROM result  LIMIT $1 OFFSET $2", paginationQuery.Take, paginationQuery.Skip)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var results []model.Result
	for rows.Next() {
		var result model.Result
		err := rows.Scan(&result.ResultId, &result.Name)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		results = append(results, result)
	}

	// count bootcamp
	var totalRows int
	row := r.db.QueryRow("SELECT COUNT(*) FROM result")
	err = row.Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	return results, common.Paginate(paginationQuery.Page, paginationQuery.Take, totalRows), nil

}
func NewResultRepository(db *sql.DB) ResultRepository {
	return &resultRepository{db: db}
}
