package repository

import (
	"database/sql"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/utils/common"
)

type FormRepository interface {
	Create(payload model.Form) error
	Get(id string) (model.Form, error)
	Delete(id string) error
	Update(payload model.Form) error
	BaseRepositoryPaging[model.Form]
}

type formRepository struct {
	db *sql.DB
}

func (f *formRepository) Create(payload model.Form) error {
	_, err := f.db.Exec("INSERT INTO form (id, form_link) VALUES ($1, $2)", payload.FormID, payload.FormLink)
	if err != nil {
		return err 
	}
	return nil
}

func (f *formRepository) Get(id string) (model.Form, error) {
	var form model.Form
	err := f.db.QueryRow(
		"SELECT id, form_link FROM form WHERE id = $1",
		id,
	).Scan(&form.FormID, &form.FormLink)
	return form, err
}

func (f *formRepository) Paging(requestPaging dto.PaginationParam) ([]model.Form, dto.Paging, error) {
	var paginationQuery dto.PaginationQuery
	paginationQuery = common.GetPaginationParams(requestPaging)
	rows, err := f.db.Query("SELECT * FROM form  LIMIT $1 OFFSET $2", paginationQuery.Take, paginationQuery.Skip)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var forms []model.Form
	for rows.Next() {
		var form model.Form
		err := rows.Scan(&form.FormID, &form.FormLink)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		forms = append(forms, form)
	}

	// count bootcamp
	var totalRows int
	row := f.db.QueryRow("SELECT COUNT(*) FROM status")
	err = row.Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	return forms, common.Paginate(paginationQuery.Page, paginationQuery.Take, totalRows), nil

}

func (f *formRepository) Update(payload model.Form) error {
	_, err := f.db.Exec(
		"UPDATE form SET form_link = $2 WHERE id = $1",
		payload.FormID, payload.FormLink,
	)
	return err
}

func (f *formRepository) Delete(id string) error {
	_, err := f.db.Exec("DELETE FROM form WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func NewFormRepository(db *sql.DB) FormRepository {
	return &formRepository{db: db}
}