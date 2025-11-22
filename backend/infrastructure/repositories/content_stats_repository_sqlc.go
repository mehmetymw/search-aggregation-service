package repositories

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/mehmetymw/search-aggregation-service/backend/db/generated"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
)

type ContentStatsRepositorySqlc struct {
	db      *sql.DB
	queries *db.Queries
}

func NewContentStatsRepository(database *sql.DB) ports.ContentStatsRepository {
	return &ContentStatsRepositorySqlc{
		db:      database,
		queries: db.New(database),
	}
}

func (r *ContentStatsRepositorySqlc) SaveOrUpdateStats(ctx context.Context, stats []entity.ContentStats) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	for _, stat := range stats {
		err := qtx.UpsertContentStats(ctx, db.UpsertContentStatsParams{
			ContentID:   stat.ContentID,
			Views:       stat.Views,
			Likes:       stat.Likes,
			DurationSec: stat.DurationSec,
			ReadingTime: stat.ReadingTime,
			Reactions:   stat.Reactions,
			Comments:    stat.Comments,
		})
		if err != nil {
			return fmt.Errorf("upsert content stats: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *ContentStatsRepositorySqlc) GetByContentIDs(ctx context.Context, contentIDs []int64) (map[int64]entity.ContentStats, error) {
	rows, err := r.queries.GetContentStatsByIDs(ctx, contentIDs)
	if err != nil {
		return nil, fmt.Errorf("get content stats by ids: %w", err)
	}

	statsMap := make(map[int64]entity.ContentStats)
	for _, row := range rows {
		statsMap[row.ContentID] = entity.ContentStats{
			ContentID:   row.ContentID,
			Views:       row.Views,
			Likes:       row.Likes,
			DurationSec: row.DurationSec,
			ReadingTime: row.ReadingTime,
			Reactions:   row.Reactions,
			Comments:    row.Comments,
			LastSyncAt:  row.LastSyncAt,
		}
	}

	return statsMap, nil
}

func (r *ContentStatsRepositorySqlc) GetByContentID(ctx context.Context, contentID int64) (*entity.ContentStats, error) {
	row, err := r.queries.GetContentStatsByID(ctx, contentID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get content stats by id: %w", err)
	}

	stats := entity.ContentStats{
		ContentID:   row.ContentID,
		Views:       row.Views,
		Likes:       row.Likes,
		DurationSec: row.DurationSec,
		ReadingTime: row.ReadingTime,
		Reactions:   row.Reactions,
		Comments:    row.Comments,
		LastSyncAt:  row.LastSyncAt,
	}
	return &stats, nil
}
