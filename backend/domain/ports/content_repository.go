package ports

import (
	"context"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
)

type Pagination struct {
	Page     int32
	PageSize int32
}

func (p Pagination) Offset() int32 {
	return (p.Page - 1) * p.PageSize
}

func (p Pagination) Limit() int32 {
	return p.PageSize
}

type SearchFilters struct {
	Query       string
	ContentType *entity.ContentType
}

type ContentRepository interface {
	SaveOrUpdateContents(ctx context.Context, contents []entity.Content) error
	SearchContents(ctx context.Context, filters SearchFilters, pagination Pagination) ([]entity.Content, int64, error)
	GetByIDs(ctx context.Context, ids []int64) ([]entity.Content, error)
	GetByID(ctx context.Context, id int64) (*entity.Content, error)
}
