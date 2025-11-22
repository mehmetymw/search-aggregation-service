package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagRepository_EnsureTags(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTagRepository(db)
	ctx := context.Background()

	tags := []string{"go", "testing", "integration"}
	
	createdTags, err := repo.EnsureTags(ctx, tags)
	assert.NoError(t, err)
	assert.Len(t, createdTags, 3)
	
	for _, tag := range createdTags {
		assert.NotEmpty(t, tag.ID)
		assert.NotEmpty(t, tag.Name)
	}
	
	// Ensure idempotency
	createdTags2, err := repo.EnsureTags(ctx, tags)
	assert.NoError(t, err)
	assert.Len(t, createdTags2, 3)
	
	// IDs should match
	assert.Equal(t, createdTags[0].ID, createdTags2[0].ID)
}
