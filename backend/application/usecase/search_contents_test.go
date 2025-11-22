package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchContentsUseCase_Execute(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockStatsRepo := new(MockContentStatsRepository)
	mockCache := new(MockCacheClient)
	mockLogger := new(MockLogger)

	scoringConfig := entity.ScoringConfig{
		VideoViewsDivisor:   1.0,
		VideoLikesDivisor:   1.0,
		VideoTypeMultiplier: 1.0,
	}
	timeProvider := func() time.Time { return time.Now() }
	scoringService := service.NewScoringService(scoringConfig, timeProvider)

	uc := NewSearchContentsUseCase(
		mockContentRepo,
		mockStatsRepo,
		mockCache,
		scoringService,
		mockLogger,
		time.Minute,
	)

	ctx := context.Background()
	req := SearchContentsRequest{
		Query:    "test",
		Page:     1,
		PageSize: 10,
		Sort:     SortScoreDesc,
	}

	t.Run("Cache Hit", func(t *testing.T) {
		cachedResult := SearchResult{Total: 100}
		mockCache.On("Get", ctx, mock.AnythingOfType("string"), mock.Anything).Return(true, nil).Run(func(args mock.Arguments) {
			dest := args.Get(2).(*SearchResult)
			*dest = cachedResult
		}).Once()

		res, err := uc.Execute(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), res.Total)
		mockCache.AssertExpectations(t)
	})

	t.Run("Cache Miss - Success", func(t *testing.T) {
		mockCache.On("Get", ctx, mock.AnythingOfType("string"), mock.Anything).Return(false, nil)

		contents := []entity.Content{
			{ID: 1, ContentType: entity.ContentTypeVideo, Title: "Video 1"},
			{ID: 2, ContentType: entity.ContentTypeVideo, Title: "Video 2"},
		}
		mockContentRepo.On("SearchContents", ctx, mock.Anything, mock.Anything).Return(contents, int64(2), nil)

		stats := map[int64]entity.ContentStats{
			1: {ContentID: 1, Views: 100, Likes: 10},
			2: {ContentID: 2, Views: 200, Likes: 20},
		}
		mockStatsRepo.On("GetByContentIDs", ctx, []int64{1, 2}).Return(stats, nil)

		mockCache.On("Set", ctx, mock.AnythingOfType("string"), mock.Anything, time.Minute).Return(nil)

		res, err := uc.Execute(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), res.Total)
		assert.Len(t, res.Items, 2)
		
		// Verify sorting (Score Desc)
		// Video 2 should have higher score (200 views > 100 views)
		assert.Equal(t, int64(2), res.Items[0].Content.ID)
		assert.Equal(t, int64(1), res.Items[1].Content.ID)
	})
}
