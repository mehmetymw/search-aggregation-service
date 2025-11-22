package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/application/usecase"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/service"
	contentpb "github.com/mehmetymw/search-aggregation-service/backend/proto/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContentServiceServer_SearchContents(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockStatsRepo := new(MockContentStatsRepository)
	mockCache := new(MockCacheClient)
	mockLogger := new(MockLogger)
	mockMetadataRepo := new(MockMetadataRepository)

	scoringConfig := entity.ScoringConfig{
		VideoViewsDivisor: 1.0,
	}
	timeProvider := func() time.Time { return time.Now() }
	scoringService := service.NewScoringService(scoringConfig, timeProvider)

	searchUC := usecase.NewSearchContentsUseCase(
		mockContentRepo,
		mockStatsRepo,
		mockCache,
		scoringService,
		mockLogger,
		time.Minute,
	)

	getByIDUC := usecase.NewGetContentByIDUseCase(
		mockContentRepo,
		mockStatsRepo,
		scoringService,
	)

	appConfig := entity.AppConfig{
		Pagination: entity.PaginationConfig{
			DefaultPage:     1,
			DefaultPageSize: 10,
			MaxPageSize:     50,
		},
	}

	server := NewContentServiceServer(
		searchUC,
		getByIDUC,
		mockMetadataRepo,
		appConfig,
		mockLogger,
	)

	ctx := context.Background()
	req := &contentpb.SearchRequest{
		Query:    "test",
		Page:     1,
		PageSize: 10,
		Sort:     "score_desc",
	}

	t.Run("Success", func(t *testing.T) {
		// Mock Cache Miss
		mockCache.On("Get", ctx, mock.AnythingOfType("string"), mock.Anything).Return(false, nil)

		// Mock Repo Search
		contents := []entity.Content{
			{ID: 1, Title: "Test Video", ContentType: entity.ContentTypeVideo},
		}
		mockContentRepo.On("SearchContents", ctx, mock.Anything, mock.Anything).Return(contents, int64(1), nil)

		// Mock Repo Stats
		stats := map[int64]entity.ContentStats{
			1: {ContentID: 1, Views: 100},
		}
		mockStatsRepo.On("GetByContentIDs", ctx, []int64{1}).Return(stats, nil)

		// Mock Cache Set
		mockCache.On("Set", ctx, mock.AnythingOfType("string"), mock.Anything, time.Minute).Return(nil)

		resp, err := server.SearchContents(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(1), resp.Total)
		assert.Len(t, resp.Items, 1)
		assert.Equal(t, "Test Video", resp.Items[0].Title)
	})
}

func TestContentServiceServer_GetContent(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockStatsRepo := new(MockContentStatsRepository)
	mockLogger := new(MockLogger)
	
	// We need to setup the server similarly...
	// For brevity, I'll just setup what's needed for GetContent
	
	scoringConfig := entity.ScoringConfig{}
	timeProvider := func() time.Time { return time.Now() }
	scoringService := service.NewScoringService(scoringConfig, timeProvider)
	
	getByIDUC := usecase.NewGetContentByIDUseCase(
		mockContentRepo,
		mockStatsRepo,
		scoringService,
	)
	
	server := &ContentServiceServer{
		getByIDUseCase: getByIDUC,
		logger:         mockLogger,
	}
	
	ctx := context.Background()
	req := &contentpb.GetContentRequest{Id: 1}
	
	t.Run("Found", func(t *testing.T) {
		content := &entity.Content{ID: 1, Title: "Found"}
		mockContentRepo.On("GetByID", ctx, int64(1)).Return(content, nil)
		
		stats := &entity.ContentStats{ContentID: 1}
		mockStatsRepo.On("GetByContentID", ctx, int64(1)).Return(stats, nil)
		
		resp, err := server.GetContent(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Content)
		assert.Equal(t, "Found", resp.Content.Title)
	})
}
