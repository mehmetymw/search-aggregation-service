package repositories

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/mehmetymw/search-aggregation-service/backend/db/generated"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
)

type ContentRepositorySqlc struct {
	db      *sql.DB
	queries *db.Queries
}

func NewContentRepository(database *sql.DB) ports.ContentRepository {
	return &ContentRepositorySqlc{
		db:      database,
		queries: db.New(database),
	}
}

func (r *ContentRepositorySqlc) SaveOrUpdateContents(ctx context.Context, contents []entity.Content) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	for _, content := range contents {
		_, err := qtx.UpsertContent(ctx, db.UpsertContentParams{
			ProviderID:        content.ProviderID,
			ProviderContentID: content.ProviderContentID,
			Title:             content.Title,
			ContentType:       string(content.ContentType),
			PublishedAt:       content.PublishedAt,
			IsActive:          content.IsActive,
		})
		if err != nil {
			return fmt.Errorf("upsert content: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *ContentRepositorySqlc) SearchContents(ctx context.Context, filters ports.SearchFilters, pagination ports.Pagination) ([]entity.Content, int64, error) {
	var queryParam sql.NullString
	if filters.Query != "" {
		queryParam = sql.NullString{String: filters.Query, Valid: true}
	}

	var contentTypeParam sql.NullString
	if filters.ContentType != nil {
		contentTypeParam = sql.NullString{
			String: string(*filters.ContentType),
			Valid:  true,
		}
	}

	rows, err := r.queries.SearchContents(ctx, db.SearchContentsParams{
		Query:       queryParam,
		ContentType: contentTypeParam,
		LimitCount:  pagination.Limit(),
		OffsetCount: pagination.Offset(),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("search contents: %w", err)
	}

	contents := make([]entity.Content, 0, len(rows))
	for _, row := range rows {
		contents = append(contents, dbRowToContent(row))
	}

	count, err := r.queries.CountContents(ctx, db.CountContentsParams{
		Query:       queryParam,
		ContentType: contentTypeParam,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("count contents: %w", err)
	}

	return contents, count, nil
}

func (r *ContentRepositorySqlc) GetByIDs(ctx context.Context, ids []int64) ([]entity.Content, error) {
	rows, err := r.queries.GetContentsByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get contents by ids: %w", err)
	}

	contents := make([]entity.Content, 0, len(rows))
	for _, row := range rows {
		contents = append(contents, dbRowToContent(row))
	}

	return contents, nil
}

func (r *ContentRepositorySqlc) GetByID(ctx context.Context, id int64) (*entity.Content, error) {
	row, err := r.queries.GetContentByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get content by id: %w", err)
	}

	content := dbRowToContent(row)
	return &content, nil
}

func dbRowToContent(row db.Content) entity.Content {
	contentTypeStr := row.ContentType
	return entity.Content{
		ID:                row.ID,
		ProviderID:        row.ProviderID,
		ProviderContentID: row.ProviderContentID,
		Title:             row.Title,
		ContentType:       entity.ContentType(contentTypeStr),
		PublishedAt:       row.PublishedAt,
		IsActive:          row.IsActive,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}
