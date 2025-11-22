package usecase

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/service"
	loggerPkg "github.com/mehmetymw/search-aggregation-service/backend/infrastructure/logger"
)

type ContentWithScore struct {
	Content entity.Content
	Stats   entity.ContentStats
	Score   entity.ScoreComponents
}

type SearchResult struct {
	Items    []ContentWithScore
	Page     int32
	PageSize int32
	Total    int64
}

type SortOption string

const (
	SortScoreDesc   SortOption = "score_desc"
	SortScoreAsc    SortOption = "score_asc"
	SortRecencyDesc SortOption = "recency_desc"
)

type SearchContentsRequest struct {
	Query       string
	ContentType *entity.ContentType
	Sort        SortOption
	Page        int32
	PageSize    int32
}

type SearchContentsUseCase struct {
	contentRepo      ports.ContentRepository
	contentStatsRepo ports.ContentStatsRepository
	cacheClient      ports.CacheClient
	scoringService   *service.ScoringService
	logger           ports.Logger
	cacheTTL         time.Duration
}

func NewSearchContentsUseCase(
	contentRepo ports.ContentRepository,
	contentStatsRepo ports.ContentStatsRepository,
	cacheClient ports.CacheClient,
	scoringService *service.ScoringService,
	logger ports.Logger,
	cacheTTL time.Duration,
) *SearchContentsUseCase {
	return &SearchContentsUseCase{
		contentRepo:      contentRepo,
		contentStatsRepo: contentStatsRepo,
		cacheClient:      cacheClient,
		scoringService:   scoringService,
		logger:           logger,
		cacheTTL:         cacheTTL,
	}
}

func (uc *SearchContentsUseCase) Execute(ctx context.Context, req SearchContentsRequest) (*SearchResult, error) {
	cacheKey := uc.buildCacheKey(req)
	
	var cachedResult SearchResult
	found, err := uc.cacheClient.Get(ctx, cacheKey, &cachedResult)
	if err == nil && found {
		return &cachedResult, nil
	}

	filters := ports.SearchFilters{
		Query:       req.Query,
		ContentType: req.ContentType,
	}

	pagination := ports.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	contents, total, err := uc.contentRepo.SearchContents(ctx, filters, pagination)
	if err != nil {
		return nil, fmt.Errorf("search contents: %w", err)
	}

	if len(contents) == 0 {
		result := &SearchResult{
			Items:    []ContentWithScore{},
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    0,
		}
		return result, nil
	}

	contentIDs := make([]int64, len(contents))
	for i, content := range contents {
		contentIDs[i] = content.ID
	}

	statsMap, err := uc.contentStatsRepo.GetByContentIDs(ctx, contentIDs)
	if err != nil {
		return nil, fmt.Errorf("get content stats: %w", err)
	}

	items := make([]ContentWithScore, 0, len(contents))
	for _, content := range contents {
		stats, ok := statsMap[content.ID]
		if !ok {
			stats = entity.ContentStats{ContentID: content.ID}
		}

		score := uc.scoringService.Calculate(content, stats)

		items = append(items, ContentWithScore{
			Content: content,
			Stats:   stats,
			Score:   score,
		})
	}

	uc.sortItems(items, req.Sort)

	result := &SearchResult{
		Items:    items,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}

	if err := uc.cacheClient.Set(ctx, cacheKey, result, uc.cacheTTL); err != nil {
		uc.logger.Warn("failed to cache search result", loggerPkg.String("error", err.Error()))
	}

	return result, nil
}

func (uc *SearchContentsUseCase) sortItems(items []ContentWithScore, sortOption SortOption) {
	switch sortOption {
	case SortScoreDesc:
		sort.Slice(items, func(i, j int) bool {
			return items[i].Score.FinalScore > items[j].Score.FinalScore
		})
	case SortScoreAsc:
		sort.Slice(items, func(i, j int) bool {
			return items[i].Score.FinalScore < items[j].Score.FinalScore
		})
	case SortRecencyDesc:
		sort.Slice(items, func(i, j int) bool {
			return items[i].Content.PublishedAt.After(items[j].Content.PublishedAt)
		})
	default:
		sort.Slice(items, func(i, j int) bool {
			return items[i].Score.FinalScore > items[j].Score.FinalScore
		})
	}
}

func (uc *SearchContentsUseCase) buildCacheKey(req SearchContentsRequest) string {
	typeStr := "all"
	if req.ContentType != nil {
		typeStr = string(*req.ContentType)
	}
	return fmt.Sprintf("search:%s:%s:%s:%d:%d",
		req.Query, typeStr, req.Sort, req.Page, req.PageSize)
}
