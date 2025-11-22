package ports

import (
	"context"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
)

type TagRepository interface {
	EnsureTags(ctx context.Context, tagNames []string) ([]entity.Tag, error)
	AssignToContent(ctx context.Context, contentID int64, tagIDs []int64) error
	GetByContentID(ctx context.Context, contentID int64) ([]entity.Tag, error)
}
