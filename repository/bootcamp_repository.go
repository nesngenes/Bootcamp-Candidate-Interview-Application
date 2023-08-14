package repository

import (
	"database/sql"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/utils/common"
)

type BootcampRepository interface {
	BaseRepository[model.Bootcamp]
	BaseRepositoryPaging[model.Bootcamp]
	GetByName(name string) (model.Bootcamp, error)
	GetByID(id string) (model.Bootcamp, error)
}

type bootcampRepository struct {
	db *sql.DB
}

func (b *bootcampRepository) Create(payload model.Bootcamp) error {
	_, err := b.db.Exec("INSERT INTO bootcamp (id, name,start_date,end_date,location) VALUES ($1, $2, $3, $4, $5)", payload.BootcampId, payload.Name, payload.StartDate, payload.EndDate, payload.Location)
	if err != nil {
		return err
	}
	return nil
}
func (b *bootcampRepository) GetByName(name string) (model.Bootcamp, error) {
	var bootcamp model.Bootcamp
	err := b.db.QueryRow("SELECT * FROM bootcamp WHERE name ILIKE $1", "%"+name+"%").Scan(&bootcamp.BootcampId, &bootcamp.Name, &bootcamp.StartDate, &bootcamp.EndDate, &bootcamp.Location)
	if err != nil {
		return model.Bootcamp{}, err
	}
	return bootcamp, nil
}

func (b *bootcampRepository) GetByID(id string) (model.Bootcamp, error) {
	var bootcamp model.Bootcamp
	err := b.db.QueryRow("SELECT * FROM bootcamp WHERE id = $1", id).
		Scan(&bootcamp.BootcampId, &bootcamp.Name, &bootcamp.StartDate, &bootcamp.EndDate, &bootcamp.Location)
	return bootcamp, err
}

func (b *bootcampRepository) List() ([]model.Bootcamp, error) {
	rows, err := b.db.Query("SELECT * FROM bootcamp")
	if err != nil {
		return nil, err
	}
	var bootcamps []model.Bootcamp
	for rows.Next() {
		var bootcamp model.Bootcamp
		err := rows.Scan(&bootcamp.BootcampId, &bootcamp.Name, &bootcamp.StartDate, &bootcamp.EndDate, &bootcamp.Location)
		if err != nil {
			return nil, err
		}
		bootcamps = append(bootcamps, bootcamp)
	}
	return bootcamps, nil
}
func (b *bootcampRepository) Get(id string) (model.Bootcamp, error) {
	var bootcamp model.Bootcamp
	row := b.db.QueryRow("SELECT * FROM bootcamp WHERE id= $1", id)
	err := row.Scan(&bootcamp.BootcampId, &bootcamp.Name, &bootcamp.StartDate, &bootcamp.EndDate, &bootcamp.Location)
	if err != nil {
		return model.Bootcamp{}, err
	}
	return bootcamp, nil

}
func (b *bootcampRepository) Update(payload model.Bootcamp) error {
	_, err := b.db.Exec("UPDATE bootcamp SET name = $2, start_date = $3, end_date = $4,location= $5  WHERE id = $1", payload.BootcampId, payload.Name, payload.StartDate, payload.EndDate, payload.Location)
	if err != nil {
		return err
	}
	return nil

}
func (b *bootcampRepository) Delete(id string) error {
	_, err := b.db.Exec("DELETE FROM bootcamp WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
func (b *bootcampRepository) Paging(requestPaging dto.PaginationParam) ([]model.Bootcamp, dto.Paging, error) {
	var paginationQuery dto.PaginationQuery
	paginationQuery = common.GetPaginationParams(requestPaging)
	rows, err := b.db.Query("SELECT * FROM bootcamp  LIMIT $1 OFFSET $2", paginationQuery.Take, paginationQuery.Skip)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var bootcamps []model.Bootcamp
	for rows.Next() {
		var bootcamp model.Bootcamp
		err := rows.Scan(&bootcamp.BootcampId, &bootcamp.Name, &bootcamp.StartDate, &bootcamp.EndDate, &bootcamp.Location)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		bootcamps = append(bootcamps, bootcamp)
	}

	// count bootcamp
	var totalRows int
	row := b.db.QueryRow("SELECT COUNT(*) FROM bootcamp")
	err = row.Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	return bootcamps, common.Paginate(paginationQuery.Page, paginationQuery.Take, totalRows), nil
}

func NewBootcampRepository(db *sql.DB) BootcampRepository {
	return &bootcampRepository{db: db}
}
