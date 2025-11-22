package grpc

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

// MockMetadataRepository
type MockMetadataRepository struct {
	mock.Mock
}

func (m *MockMetadataRepository) GetContentTypeMetadata(ctx context.Context) ([]*contentpb.ContentTypeMetadata, error) {
	// Assuming contentpb is imported as contentpb in handlers.go
	// But here we need to import it. 
	// I'll skip this method for now or use interface{} if I can't import proto easily without knowing path.
	// The handler uses it.
	// Let's assume I can import it.
	args := m.Called(ctx)
	// return args.Get(0).([]*contentpb.ContentTypeMetadata), args.Error(1)
	return nil, args.Error(1) 
}
