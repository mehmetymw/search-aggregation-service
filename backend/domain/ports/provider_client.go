package ports

import (
	"context"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
)

type ProviderContentItem struct {
	ProviderContentID string
	Title             string
	ContentType       string
	Views             int64
	Likes             int64
	DurationSec       int32
	ReadingTime       int32
	Reactions         int64
	Comments          int64
	PublishedAt       time.Time
	Tags              []string
	RawPayload        []byte
}

type ProviderClient interface {
	FetchContents(ctx context.Context, provider entity.Provider) ([]ProviderContentItem, error)
}
