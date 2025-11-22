package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
	"github.com/stretchr/testify/assert"
)

func TestContentRepository_SaveAndSearch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewContentRepository(db)
	ctx := context.Background()

	// Cleanup
	// In a real scenario, we might truncate tables. 
	// For now, we rely on unique constraints or just adding new data.

	content := entity.Content{
		ProviderID:        1,
		ProviderContentID: "test-content-1",
		Title:             "Integration Test Content",
		ContentType:       entity.ContentTypeVideo,
		PublishedAt:       time.Now().UTC(),
		IsActive:          true,
	}

	// Test Save
	err := repo.SaveOrUpdateContents(ctx, []entity.Content{content})
	assert.NoError(t, err)

	// Test Search
	filters := ports.SearchFilters{
		Query: "Integration Test",
	}
	pagination := ports.Pagination{
		Page:     1,
		PageSize: 10,
	}

	results, count, err := repo.SearchContents(ctx, filters, pagination)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))
	assert.NotEmpty(t, results)
	
	found := false
	for _, c := range results {
		if c.ProviderContentID == content.ProviderContentID {
			found = true
			assert.Equal(t, content.Title, c.Title)
			break
		}
	}
	assert.True(t, found, "saved content not found in search results")
}
