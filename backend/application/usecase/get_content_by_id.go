package usecase

import (
	"context"
	"fmt"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/service"
)

type GetContentByIDRequest struct {
	ID int64
}

type GetContentByIDUseCase struct {
	contentRepo      ports.ContentRepository
	contentStatsRepo ports.ContentStatsRepository
	scoringService   *service.ScoringService
}

func NewGetContentByIDUseCase(
	contentRepo ports.ContentRepository,
	contentStatsRepo ports.ContentStatsRepository,
	scoringService *service.ScoringService,
) *GetContentByIDUseCase {
	return &GetContentByIDUseCase{
		contentRepo:      contentRepo,
		contentStatsRepo: contentStatsRepo,
		scoringService:   scoringService,
	}
}

func (uc *GetContentByIDUseCase) Execute(ctx context.Context, req GetContentByIDRequest) (*ContentWithScore, error) {
	content, err := uc.contentRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("get content: %w", err)
	}
	if content == nil {
		return nil, nil
	}

	stats, err := uc.contentStatsRepo.GetByContentID(ctx, content.ID)
	if err != nil {
		return nil, fmt.Errorf("get content stats: %w", err)
	}
	if stats == nil {
		stats = &entity.ContentStats{ContentID: content.ID}
	}

	score := uc.scoringService.Calculate(*content, *stats)

	return &ContentWithScore{
		Content: *content,
		Stats:   *stats,
		Score:   score,
	}, nil
}
