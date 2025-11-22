package usecase

import (
	"context"
	"testing"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSyncProviderContentsUseCase_ExecuteForProvider(t *testing.T) {
	mockProviderRepo := new(MockProviderRepository)
	mockContentRepo := new(MockContentRepository)
	mockStatsRepo := new(MockContentStatsRepository)
	mockTagRepo := new(MockTagRepository)
	mockJsonClient := new(MockProviderClient)
	mockXmlClient := new(MockProviderClient)
	mockLogger := new(MockLogger)
	tagNormalizer := service.NewTagNormalizer()

	uc := NewSyncProviderContentsUseCase(
		mockProviderRepo,
		mockContentRepo,
		mockStatsRepo,
		mockTagRepo,
		mockJsonClient,
		mockXmlClient,
		tagNormalizer,
		mockLogger,
	)

	ctx := context.Background()
	provider := entity.Provider{
		ID:     1,
		Code:   "provider1",
		Format: entity.ProviderFormatJSON,
	}

	t.Run("Success", func(t *testing.T) {
		items := []ports.ProviderContentItem{
			{
				ProviderContentID: "p1",
				Title:             "Title 1",
				Tags:              []string{"Tag1", "Tag2"},
			},
		}
		mockJsonClient.On("FetchContents", ctx, provider).Return(items, nil)

		mockContentRepo.On("SaveOrUpdateContents", ctx, mock.Anything).Return(nil)
		
		// Mock searching back the saved contents to get IDs
		savedContents := []entity.Content{
			{ID: 101, ProviderID: 1, ProviderContentID: "p1"},
		}
		mockContentRepo.On("SearchContents", ctx, mock.Anything, mock.Anything).Return(savedContents, int64(1), nil)

		mockStatsRepo.On("SaveOrUpdateStats", ctx, mock.Anything).Return(nil)

		mockTagRepo.On("EnsureTags", ctx, mock.Anything).Return([]entity.Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}, nil)
		mockTagRepo.On("AssignToContent", ctx, int64(101), mock.Anything).Return(nil)

		mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()

		err := uc.ExecuteForProvider(ctx, provider)
		assert.NoError(t, err)
	})
}
