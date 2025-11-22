package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/service"
	"github.com/stretchr/testify/assert"
)

func TestGetContentByIDUseCase_Execute(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockStatsRepo := new(MockContentStatsRepository)

	scoringConfig := entity.ScoringConfig{}
	timeProvider := func() time.Time { return time.Now() }
	scoringService := service.NewScoringService(scoringConfig, timeProvider)

	uc := NewGetContentByIDUseCase(
		mockContentRepo,
		mockStatsRepo,
		scoringService,
	)

	ctx := context.Background()

	t.Run("Found", func(t *testing.T) {
		content := &entity.Content{ID: 1, Title: "Test"}
		mockContentRepo.On("GetByID", ctx, int64(1)).Return(content, nil)

		stats := &entity.ContentStats{ContentID: 1, Views: 100}
		mockStatsRepo.On("GetByContentID", ctx, int64(1)).Return(stats, nil)

		res, err := uc.Execute(ctx, GetContentByIDRequest{ID: 1})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, int64(1), res.Content.ID)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockContentRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

		res, err := uc.Execute(ctx, GetContentByIDRequest{ID: 999})
		assert.NoError(t, err)
		assert.Nil(t, res)
	})
}
