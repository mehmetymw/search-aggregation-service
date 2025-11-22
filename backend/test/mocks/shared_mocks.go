package mocks

import (
	"context"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
	contentpb "github.com/mehmetymw/search-aggregation-service/backend/proto/gen"
	"github.com/stretchr/testify/mock"
)

// MockContentRepository
type MockContentRepository struct {
	mock.Mock
}

func (m *MockContentRepository) SaveOrUpdateContents(ctx context.Context, contents []entity.Content) error {
	args := m.Called(ctx, contents)
	return args.Error(0)
}

func (m *MockContentRepository) SearchContents(ctx context.Context, filters ports.SearchFilters, pagination ports.Pagination) ([]entity.Content, int64, error) {
	args := m.Called(ctx, filters, pagination)
	return args.Get(0).([]entity.Content), args.Get(1).(int64), args.Error(2)
}

func (m *MockContentRepository) GetByIDs(ctx context.Context, ids []int64) ([]entity.Content, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]entity.Content), args.Error(1)
}

func (m *MockContentRepository) GetByID(ctx context.Context, id int64) (*entity.Content, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Content), args.Error(1)
}

// MockContentStatsRepository
type MockContentStatsRepository struct {
	mock.Mock
}

func (m *MockContentStatsRepository) SaveOrUpdateStats(ctx context.Context, stats []entity.ContentStats) error {
	args := m.Called(ctx, stats)
	return args.Error(0)
}

func (m *MockContentStatsRepository) GetByContentIDs(ctx context.Context, contentIDs []int64) (map[int64]entity.ContentStats, error) {
	args := m.Called(ctx, contentIDs)
	return args.Get(0).(map[int64]entity.ContentStats), args.Error(1)
}

func (m *MockContentStatsRepository) GetByContentID(ctx context.Context, contentID int64) (*entity.ContentStats, error) {
	args := m.Called(ctx, contentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ContentStats), args.Error(1)
}

// MockCacheClient
type MockCacheClient struct {
	mock.Mock
}

func (m *MockCacheClient) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	args := m.Called(ctx, key, dest)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheClient) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// MockLogger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...ports.Field) {
	args := make([]interface{}, len(fields)+1)
	args[0] = msg
	for i, f := range fields {
		args[i+1] = f
	}
	m.Called(args...)
}

func (m *MockLogger) Error(msg string, fields ...ports.Field) {
	args := make([]interface{}, len(fields)+1)
	args[0] = msg
	for i, f := range fields {
		args[i+1] = f
	}
	m.Called(args...)
}

func (m *MockLogger) Warn(msg string, fields ...ports.Field) {
	args := make([]interface{}, len(fields)+1)
	args[0] = msg
	for i, f := range fields {
		args[i+1] = f
	}
	m.Called(args...)
}

func (m *MockLogger) Debug(msg string, fields ...ports.Field) {
	args := make([]interface{}, len(fields)+1)
	args[0] = msg
	for i, f := range fields {
		args[i+1] = f
	}
	m.Called(args...)
}

func (m *MockLogger) With(fields ...ports.Field) ports.Logger {
	args := make([]interface{}, len(fields))
	for i, f := range fields {
		args[i] = f
	}
	m.Called(args...)
	return m
}

// MockProviderRepository
type MockProviderRepository struct {
	mock.Mock
}

func (m *MockProviderRepository) GetAllEnabled(ctx context.Context) ([]entity.Provider, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Provider), args.Error(1)
}

func (m *MockProviderRepository) GetByCode(ctx context.Context, code string) (*entity.Provider, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Provider), args.Error(1)
}

func (m *MockProviderRepository) GetByID(ctx context.Context, id int64) (*entity.Provider, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Provider), args.Error(1)
}

func (m *MockProviderRepository) UpsertProvider(ctx context.Context, provider entity.Provider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

// MockTagRepository
type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) EnsureTags(ctx context.Context, tags []string) ([]entity.Tag, error) {
	args := m.Called(ctx, tags)
	return args.Get(0).([]entity.Tag), args.Error(1)
}

func (m *MockTagRepository) AssignToContent(ctx context.Context, contentID int64, tagIDs []int64) error {
	args := m.Called(ctx, contentID, tagIDs)
	return args.Error(0)
}

func (m *MockTagRepository) GetByContentID(ctx context.Context, contentID int64) ([]entity.Tag, error) {
	args := m.Called(ctx, contentID)
	return args.Get(0).([]entity.Tag), args.Error(1)
}

// MockProviderClient
type MockProviderClient struct {
	mock.Mock
}

func (m *MockProviderClient) FetchContents(ctx context.Context, provider entity.Provider) ([]ports.ProviderContentItem, error) {
	args := m.Called(ctx, provider)
	return args.Get(0).([]ports.ProviderContentItem), args.Error(1)
}

// MockMetadataRepository
type MockMetadataRepository struct {
	mock.Mock
}

func (m *MockMetadataRepository) GetContentTypeMetadata(ctx context.Context) ([]*contentpb.ContentTypeMetadata, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*contentpb.ContentTypeMetadata), args.Error(1)
}
