package usecase

import (
	"context"
	"fmt"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/service"
	loggerPkg "github.com/mehmetymw/search-aggregation-service/backend/infrastructure/logger"
)

type SyncProviderContentsUseCase struct {
	providerRepo     ports.ProviderRepository
	contentRepo      ports.ContentRepository
	contentStatsRepo ports.ContentStatsRepository
	tagRepo          ports.TagRepository
	providerClients  map[string]ports.ProviderClient
	tagNormalizer    *service.TagNormalizer
	logger           ports.Logger
}

func NewSyncProviderContentsUseCase(
	providerRepo ports.ProviderRepository,
	contentRepo ports.ContentRepository,
	contentStatsRepo ports.ContentStatsRepository,
	tagRepo ports.TagRepository,
	jsonClient ports.ProviderClient,
	xmlClient ports.ProviderClient,
	tagNormalizer *service.TagNormalizer,
	logger ports.Logger,
) *SyncProviderContentsUseCase {
	return &SyncProviderContentsUseCase{
		providerRepo:     providerRepo,
		contentRepo:      contentRepo,
		contentStatsRepo: contentStatsRepo,
		tagRepo:          tagRepo,
		providerClients: map[string]ports.ProviderClient{
			entity.ProviderFormatJSON: jsonClient,
			entity.ProviderFormatXML:  xmlClient,
		},
		tagNormalizer: tagNormalizer,
		logger:        logger,
	}
}

func (uc *SyncProviderContentsUseCase) ExecuteAll(ctx context.Context) error {
	providers, err := uc.providerRepo.GetAllEnabled(ctx)
	if err != nil {
		return fmt.Errorf("get all enabled providers: %w", err)
	}

	uc.logger.Info("starting sync for all providers", loggerPkg.Int("provider_count", len(providers)))

	for _, provider := range providers {
		if err := uc.ExecuteForProvider(ctx, provider); err != nil {
			uc.logger.Error("sync failed for provider",
				loggerPkg.String("provider_code", provider.Code),
				loggerPkg.Error(err))
			continue
		}
		uc.logger.Info("synced provider successfully", loggerPkg.String("provider_code", provider.Code))
	}

	return nil
}

func (uc *SyncProviderContentsUseCase) ExecuteForProvider(ctx context.Context, provider entity.Provider) error {
	client, ok := uc.providerClients[provider.Format]
	if !ok {
		return fmt.Errorf("no client registered for provider format: %s", provider.Format)
	}

	items, err := client.FetchContents(ctx, provider)
	if err != nil {
		return fmt.Errorf("fetch contents: %w", err)
	}

	if len(items) == 0 {
		uc.logger.Info("no items fetched from provider", loggerPkg.String("provider_code", provider.Code))
		return nil
	}

	contents := make([]entity.Content, 0, len(items))
	for _, item := range items {
		contents = append(contents, entity.Content{
			ProviderID:        provider.ID,
			ProviderContentID: item.ProviderContentID,
			Title:             item.Title,
			ContentType:       entity.ContentType(item.ContentType),
			PublishedAt:       item.PublishedAt,
			IsActive:          true,
		})
	}

	if err := uc.contentRepo.SaveOrUpdateContents(ctx, contents); err != nil {
		return fmt.Errorf("save contents: %w", err)
	}

	savedContents, _, err := uc.contentRepo.SearchContents(ctx, ports.SearchFilters{}, ports.Pagination{Page: 1, PageSize: 10000})
	if err != nil {
		return fmt.Errorf("get saved contents: %w", err)
	}

	providerContentIDMap := make(map[string]int64)
	for _, content := range savedContents {
		if content.ProviderID == provider.ID {
			providerContentIDMap[content.ProviderContentID] = content.ID
		}
	}

	stats := make([]entity.ContentStats, 0, len(items))
	for _, item := range items {
		contentID, ok := providerContentIDMap[item.ProviderContentID]
		if !ok {
			continue
		}

		stats = append(stats, entity.ContentStats{
			ContentID:   contentID,
			Views:       item.Views,
			Likes:       item.Likes,
			DurationSec: item.DurationSec,
			ReadingTime: item.ReadingTime,
			Reactions:   item.Reactions,
			Comments:    item.Comments,
		})
	}

	if err := uc.contentStatsRepo.SaveOrUpdateStats(ctx, stats); err != nil {
		return fmt.Errorf("save content stats: %w", err)
	}

	for _, item := range items {
		contentID, ok := providerContentIDMap[item.ProviderContentID]
		if !ok || len(item.Tags) == 0 {
			continue
		}

		normalizedTags := uc.tagNormalizer.Normalize(item.Tags)
		tags, err := uc.tagRepo.EnsureTags(ctx, normalizedTags)
		if err != nil {
			uc.logger.Error("failed to ensure tags",
				loggerPkg.String("provider_content_id", item.ProviderContentID),
				loggerPkg.Error(err))
			continue
		}

		tagIDs := make([]int64, len(tags))
		for i, tag := range tags {
			tagIDs[i] = tag.ID
		}

		if err := uc.tagRepo.AssignToContent(ctx, contentID, tagIDs); err != nil {
			uc.logger.Error("failed to assign tags",
				loggerPkg.Int64("content_id", contentID),
				loggerPkg.Error(err))
		}
	}

	uc.logger.Info("synced provider items successfully",
		loggerPkg.String("provider_code", provider.Code),
		loggerPkg.Int("item_count", len(items)))

	return nil
}
