package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/application/usecase"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
	loggerPkg "github.com/mehmetymw/search-aggregation-service/backend/infrastructure/logger"
	contentpb "github.com/mehmetymw/search-aggregation-service/backend/proto/gen"
)

type ContentServiceServer struct {
	contentpb.UnimplementedContentServiceServer
	searchUseCase  *usecase.SearchContentsUseCase
	getByIDUseCase *usecase.GetContentByIDUseCase
	metadataRepo   ports.MetadataRepository
	logger         ports.Logger
	appConfig      entity.AppConfig
}

func NewContentServiceServer(
	searchUseCase *usecase.SearchContentsUseCase,
	getByIDUseCase *usecase.GetContentByIDUseCase,
	metadataRepo ports.MetadataRepository,
	appConfig entity.AppConfig,
	logger ports.Logger,
) *ContentServiceServer {
	return &ContentServiceServer{
		searchUseCase:  searchUseCase,
		getByIDUseCase: getByIDUseCase,
		metadataRepo:   metadataRepo,
		appConfig:      appConfig,
		logger:         logger,
	}
}

func (s *ContentServiceServer) SearchContents(ctx context.Context, req *contentpb.SearchRequest) (*contentpb.SearchResponse, error) {
	defaultPage := int32(s.appConfig.Pagination.DefaultPage)
	if defaultPage <= 0 {
		defaultPage = 1
	}

	defaultPageSize := int32(s.appConfig.Pagination.DefaultPageSize)
	if defaultPageSize <= 0 {
		defaultPageSize = 10
	}

	maxPageSize := int32(s.appConfig.Pagination.MaxPageSize)
	if maxPageSize <= 0 {
		maxPageSize = defaultPageSize
	}

	page := req.Page
	if page <= 0 {
		page = defaultPage
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	var contentType *entity.ContentType
	if req.Type != "" && req.Type != "all" {
		ct := entity.ContentType(req.Type)
		contentType = &ct
	}

	sortOption := usecase.SortOption(req.Sort)
	if sortOption == "" {
		sortOption = usecase.SortScoreDesc
	}

	useCaseReq := usecase.SearchContentsRequest{
		Query:       req.Query,
		ContentType: contentType,
		Sort:        sortOption,
		Page:        page,
		PageSize:    pageSize,
	}

	result, err := s.searchUseCase.Execute(ctx, useCaseReq)
	if err != nil {
		s.logger.Error("search failed", loggerPkg.Error(err))
		return nil, fmt.Errorf("search: %w", err)
	}

	items := make([]*contentpb.ContentItem, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, s.toProtoContentItem(item))
	}

	return &contentpb.SearchResponse{
		Items:    items,
		Page:     result.Page,
		PageSize: result.PageSize,
		Total:    result.Total,
	}, nil
}

func (s *ContentServiceServer) GetContent(ctx context.Context, req *contentpb.GetContentRequest) (*contentpb.GetContentResponse, error) {
	useCaseReq := usecase.GetContentByIDRequest{
		ID: req.Id,
	}

	result, err := s.getByIDUseCase.Execute(ctx, useCaseReq)
	if err != nil {
		s.logger.Error("get content failed", loggerPkg.Int64("id", req.Id), loggerPkg.Error(err))
		return nil, fmt.Errorf("get content: %w", err)
	}

	if result == nil {
		return &contentpb.GetContentResponse{}, nil
	}

	return &contentpb.GetContentResponse{
		Content: s.toProtoContentItem(*result),
	}, nil
}

func (s *ContentServiceServer) GetMetadata(ctx context.Context, req *contentpb.GetMetadataRequest) (*contentpb.GetMetadataResponse, error) {
	contentTypes, err := s.metadataRepo.GetContentTypeMetadata(ctx)
	if err != nil {
		s.logger.Error("failed to get content type metadata", loggerPkg.Error(err))
		return nil, fmt.Errorf("get content type metadata: %w", err)
	}

	return &contentpb.GetMetadataResponse{
		ContentTypes: contentTypes,
		SortOptions:  sortOptionMetadata,
		Pagination: &contentpb.PaginationMetadata{
			DefaultPageSize: int32(s.appConfig.Pagination.DefaultPageSize),
			MaxPageSize:     int32(s.appConfig.Pagination.MaxPageSize),
		},
	}, nil
}

func (s *ContentServiceServer) toProtoContentItem(item usecase.ContentWithScore) *contentpb.ContentItem {
	return &contentpb.ContentItem{
		Id:           item.Content.ID,
		Title:        item.Content.Title,
		ContentType:  string(item.Content.ContentType),
		Score:        item.Score.FinalScore,
		PublishedAt:  item.Content.PublishedAt.Format(time.RFC3339),
		ProviderName: fmt.Sprintf("provider-%d", item.Content.ProviderID),
	}
}
