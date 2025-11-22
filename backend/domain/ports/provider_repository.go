package ports

import (
	"context"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
)

type ProviderRepository interface {
	GetAllEnabled(ctx context.Context) ([]entity.Provider, error)
	GetByCode(ctx context.Context, code string) (*entity.Provider, error)
	GetByID(ctx context.Context, id int64) (*entity.Provider, error)
	UpsertProvider(ctx context.Context, provider entity.Provider) error
}
