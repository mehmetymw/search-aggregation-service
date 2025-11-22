package ports

import (
	"context"

	contentpb "github.com/mehmetymw/search-aggregation-service/backend/proto/gen"
)

type MetadataRepository interface {
	GetContentTypeMetadata(ctx context.Context) ([]*contentpb.ContentTypeMetadata, error)
}
