package repositories

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/mehmetymw/search-aggregation-service/backend/db/generated"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
)

type TagRepositorySqlc struct {
	db      *sql.DB
	queries *db.Queries
}

func NewTagRepository(database *sql.DB) ports.TagRepository {
	return &TagRepositorySqlc{
		db:      database,
		queries: db.New(database),
	}
}

func (r *TagRepositorySqlc) EnsureTags(ctx context.Context, tagNames []string) ([]entity.Tag, error) {
	tags := make([]entity.Tag, 0, len(tagNames))

	for _, name := range tagNames {
		row, err := r.queries.EnsureTag(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("ensure tag %s: %w", name, err)
		}
		tags = append(tags, entity.Tag{
			ID:   row.ID,
			Name: row.Name,
		})
	}

	return tags, nil
}

func (r *TagRepositorySqlc) AssignToContent(ctx context.Context, contentID int64, tagIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	if err := qtx.RemoveContentTags(ctx, contentID); err != nil {
		return fmt.Errorf("remove existing tags: %w", err)
	}

	for _, tagID := range tagIDs {
		err := qtx.AssignTagToContent(ctx, db.AssignTagToContentParams{
			ContentID: contentID,
			TagID:     tagID,
		})
		if err != nil {
			return fmt.Errorf("assign tag: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *TagRepositorySqlc) GetByContentID(ctx context.Context, contentID int64) ([]entity.Tag, error) {
	rows, err := r.queries.GetTagsByContentID(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("get tags by content id: %w", err)
	}

	tags := make([]entity.Tag, 0, len(rows))
	for _, row := range rows {
		tags = append(tags, entity.Tag{
			ID:   row.ID,
			Name: row.Name,
		})
	}

	return tags, nil
}
