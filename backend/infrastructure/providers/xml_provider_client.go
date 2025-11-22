package providers

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
)

type XmlProviderClient struct {
	client *http.Client
}

type xmlResponse struct {
	XMLName xml.Name `xml:"feed"`
	Items   struct {
		ItemList []xmlItem `xml:"item"`
	} `xml:"items"`
}

type xmlItem struct {
	ID       string `xml:"id"`
	Headline string `xml:"headline"`
	Type     string `xml:"type"`
	Stats    struct {
		Views       int64  `xml:"views"`
		Likes       int64  `xml:"likes"`
		Duration    string `xml:"duration"`
		ReadingTime int32  `xml:"reading_time"`
		Reactions   int64  `xml:"reactions"`
		Comments    int64  `xml:"comments"`
	} `xml:"stats"`
	PublicationDate string `xml:"publication_date"`
	Categories      struct {
		CategoryList []string `xml:"category"`
	} `xml:"categories"`
}

func parseXMLDuration(durationStr string) int32 {
	if durationStr == "" {
		return 0
	}

	var minutes, seconds int32
	if n, _ := fmt.Sscanf(durationStr, "%d:%d", &minutes, &seconds); n == 2 {
		return (minutes * 60) + seconds
	}

	var totalSeconds int32
	fmt.Sscanf(durationStr, "%d", &totalSeconds)
	return totalSeconds
}

func NewXmlProviderClient() ports.ProviderClient {
	return &XmlProviderClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *XmlProviderClient) FetchContents(ctx context.Context, provider entity.Provider) ([]ports.ProviderContentItem, error) {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var data xmlResponse
	if err := xml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	items := make([]ports.ProviderContentItem, 0, len(data.Items.ItemList))
	for _, item := range data.Items.ItemList {
		pubDate, err := time.Parse("2006-01-02", item.PublicationDate)
		if err != nil {
			pubDate = time.Now()
		}

		itemPayload, _ := json.Marshal(map[string]any{
			"id":               item.ID,
			"headline":         item.Headline,
			"type":             item.Type,
			"stats":            item.Stats,
			"publication_date": item.PublicationDate,
			"categories":       item.Categories.CategoryList,
		})

		items = append(items, ports.ProviderContentItem{
			ProviderContentID: item.ID,
			Title:             item.Headline,
			ContentType:       item.Type,
			Views:             item.Stats.Views,
			Likes:             item.Stats.Likes,
			DurationSec:       parseXMLDuration(item.Stats.Duration),
			ReadingTime:       item.Stats.ReadingTime,
			Reactions:         item.Stats.Reactions,
			Comments:          item.Stats.Comments,
			PublishedAt:       pubDate,
			Tags:              item.Categories.CategoryList,
			RawPayload:        itemPayload,
		})
	}

	return items, nil
}
