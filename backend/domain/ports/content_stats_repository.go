package ports

import (
	"context"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
)

type ContentStatsRepository interface {
	SaveOrUpdateStats(ctx context.Context, stats []entity.ContentStats) error
	GetByContentIDs(ctx context.Context, contentIDs []int64) (map[int64]entity.ContentStats, error)
	GetByContentID(ctx context.Context, contentID int64) (*entity.ContentStats, error)
}
