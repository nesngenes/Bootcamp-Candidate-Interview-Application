package repository

import "interview_bootcamp/model/dto"

type BaseRepository[T any] interface {
	Create(payload T) error
	List() ([]T, error)
	Get(id string) (T, error)
	Update(payload T) error
	Delete(id string) error
}

type BaseRepositoryPaging[T any] interface {
	Paging(requestPaging dto.PaginationParam) ([]T, dto.Paging, error)
}
package repository

import "interview_bootcamp/model/dto"

// Generic Type -> T (sebenarnya bebas, cuma umumnya ini dilambangkan dengan huruf T), T (UOM, PRODUCT, CUSTOMER, EMPLOYEE)
type BaseRepository[T any] interface {
	Create(payload T) error
	List() ([]T, error)
	Get(id string) (T, error)
	Update(payload T) error
	Delete(id string) error
}

type BaseRepositoryPaging[T any] interface {
	Paging(requestPaging dto.PaginationParam) ([]T, dto.Paging, error)
}
