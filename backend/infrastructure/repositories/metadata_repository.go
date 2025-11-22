package repositories

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/mehmetymw/search-aggregation-service/backend/db/generated"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
	contentpb "github.com/mehmetymw/search-aggregation-service/backend/proto/gen"
)

type MetadataRepository struct {
	queries *db.Queries
}

func NewMetadataRepository(conn *sql.DB) ports.MetadataRepository {
	return &MetadataRepository{
		queries: db.New(conn),
	}
}

func (r *MetadataRepository) GetContentTypeMetadata(ctx context.Context) ([]*contentpb.ContentTypeMetadata, error) {
	rows, err := r.queries.GetAllContentTypeMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("get content type metadata: %w", err)
	}

	result := []*contentpb.ContentTypeMetadata{
		{
			Id:          "all",
			DisplayName: "All",
		},
	}

	for _, row := range rows {
		result = append(result, &contentpb.ContentTypeMetadata{
			Id:          row.ID,
			DisplayName: row.DisplayName,
		})
	}

	return result, nil
}
