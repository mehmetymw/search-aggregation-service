package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestJsonProviderClient_FetchContents(t *testing.T) {
	mockResponse := `{
		"contents": [
			{
				"id": "1",
				"title": "Test Video",
				"type": "video",
				"metrics": {
					"views": 100,
					"likes": 10,
					"duration": "120",
					"reading_time": 0,
					"reactions": 5,
					"comments": 2
				},
				"published_at": "2023-10-25T12:00:00Z",
				"tags": ["tech", "go"]
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := NewJsonProviderClient()
	provider := entity.Provider{
		BaseURL: server.URL,
	}

	ctx := context.Background()
	items, err := client.FetchContents(ctx, provider)

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	
	item := items[0]
	assert.Equal(t, "1", item.ProviderContentID)
	assert.Equal(t, "Test Video", item.Title)
	assert.Equal(t, "video", item.ContentType)
	assert.Equal(t, int64(100), item.Views)
	assert.Equal(t, int32(120), item.DurationSec)
	assert.Len(t, item.Tags, 2)
	
	expectedTime, _ := time.Parse(time.RFC3339, "2023-10-25T12:00:00Z")
	assert.Equal(t, expectedTime, item.PublishedAt)
}

func TestJsonProviderClient_FetchContents_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewJsonProviderClient()
	provider := entity.Provider{
		BaseURL: server.URL,
	}

	ctx := context.Background()
	items, err := client.FetchContents(ctx, provider)

	assert.Error(t, err)
	assert.Nil(t, items)
}
