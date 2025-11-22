package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
)

type JsonProviderClient struct {
	client *http.Client
}

type jsonResponse struct {
	Contents []jsonItem `json:"contents"`
}

type jsonItem struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	Metrics struct {
		Views       int64  `json:"views"`
		Likes       int64  `json:"likes"`
		Duration    string `json:"duration"`
		ReadingTime int32  `json:"reading_time"`
		Reactions   int64  `json:"reactions"`
		Comments    int64  `json:"comments"`
	} `json:"metrics"`
	PublishedAt time.Time `json:"published_at"`
	Tags        []string  `json:"tags"`
}

func parseDurationString(durationStr string) int32 {
	if durationStr == "" {
		return 0
	}
	
	var duration int32
	fmt.Sscanf(durationStr, "%d", &duration)
	return duration
}

func NewJsonProviderClient() ports.ProviderClient {
	return &JsonProviderClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *JsonProviderClient) FetchContents(ctx context.Context, provider entity.Provider) ([]ports.ProviderContentItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, provider.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	rawBody, err := json.Marshal(resp.Body)
	if err != nil {
		rawBody = []byte{}
	}

	var data jsonResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	items := make([]ports.ProviderContentItem, 0, len(data.Contents))
	for _, item := range data.Contents {
		itemPayload, _ := json.Marshal(item)
		
		items = append(items, ports.ProviderContentItem{
			ProviderContentID: item.ID,
			Title:             item.Title,
			ContentType:       item.Type,
			Views:             item.Metrics.Views,
			Likes:             item.Metrics.Likes,
			DurationSec:       parseDurationString(item.Metrics.Duration),
			ReadingTime:       item.Metrics.ReadingTime,
			Reactions:         item.Metrics.Reactions,
			Comments:          item.Metrics.Comments,
			PublishedAt:       item.PublishedAt,
			Tags:              item.Tags,
			RawPayload:        itemPayload,
		})
	}

	_ = rawBody

	return items, nil
}
