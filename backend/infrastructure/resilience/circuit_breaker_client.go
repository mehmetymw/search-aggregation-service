package resilience

import (
	"context"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
	"github.com/sony/gobreaker"
)

type CircuitBreakerProviderClient struct {
	client ports.ProviderClient
	cb     *gobreaker.CircuitBreaker
}

func NewCircuitBreakerProviderClient(client ports.ProviderClient, config entity.CircuitBreakerConfig) *CircuitBreakerProviderClient {
	settings := gobreaker.Settings{
		Name:        config.Name,
		MaxRequests: config.MaxRequests,
		Interval:    time.Duration(config.Interval) * time.Second,
		Timeout:     time.Duration(config.Timeout) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}

	return &CircuitBreakerProviderClient{
		client: client,
		cb:     gobreaker.NewCircuitBreaker(settings),
	}
}

func (c *CircuitBreakerProviderClient) FetchContents(ctx context.Context, provider entity.Provider) ([]ports.ProviderContentItem, error) {
	result, err := c.cb.Execute(func() (interface{}, error) {
		return c.client.FetchContents(ctx, provider)
	})

	if err != nil {
		return nil, err
	}

	return result.([]ports.ProviderContentItem), nil
}
