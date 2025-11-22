package service

import (
	"testing"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestScoringService_Calculate(t *testing.T) {
	// Fixed time for testing
	now := time.Date(2023, 10, 25, 12, 0, 0, 0, time.UTC)
	timeProvider := func() time.Time { return now }

	config := entity.ScoringConfig{
		VideoViewsDivisor:      100.0,
		VideoLikesDivisor:      10.0,
		TextReadingTimeDivisor: 60.0,
		TextReactionsDivisor:   5.0,
		VideoTypeMultiplier:    1.5,
		TextTypeMultiplier:     1.2,
		RecencyWeekScore:       10.0,
		RecencyMonthScore:      5.0,
		RecencyQuarterScore:    2.0,
		VideoEngagementWeight:  2.0,
		TextEngagementWeight:   1.5,
	}

	service := NewScoringService(config, timeProvider)

	tests := []struct {
		name     string
		content  entity.Content
		stats    entity.ContentStats
		expected entity.ScoreComponents
	}{
		{
			name: "Video - High Engagement & Recent",
			content: entity.Content{
				ContentType: entity.ContentTypeVideo,
				PublishedAt: now.Add(-24 * time.Hour), // 1 day ago
			},
			stats: entity.ContentStats{
				Views: 1000,
				Likes: 200,
			},
			expected: entity.ScoreComponents{
				BaseScore:       (1000.0 / 100.0) + (200.0 / 10.0), // 10 + 20 = 30
				TypeMultiplier:  1.5,
				RecencyScore:    10.0,
				EngagementScore: (200.0 / 1000.0) * 2.0, // 0.2 * 2.0 = 0.4
				// Final: (30 * 1.5) + 10 + 0.4 = 45 + 10 + 0.4 = 55.4
				FinalScore: 55.4,
			},
		},
		{
			name: "Article - Moderate Engagement & Old",
			content: entity.Content{
				ContentType: entity.ContentTypeArticle,
				PublishedAt: now.Add(-40 * 24 * time.Hour), // 40 days ago (Quarter)
			},
			stats: entity.ContentStats{
				ReadingTime: 300, // 5 mins
				Reactions:   20,
			},
			expected: entity.ScoreComponents{
				BaseScore:       (300.0 / 60.0) + (20.0 / 5.0), // 5 + 4 = 9
				TypeMultiplier:  1.2,
				RecencyScore:    2.0, // Quarter score
				EngagementScore: (20.0 / 300.0) * 1.5, // 0.0666... * 1.5 = 0.1
				// Final: (9 * 1.2) + 2 + 0.1 = 10.8 + 2 + 0.1 = 12.9
				FinalScore: 12.9,
			},
		},
		{
			name: "Video - Zero Stats",
			content: entity.Content{
				ContentType: entity.ContentTypeVideo,
				PublishedAt: now.Add(-2 * 24 * time.Hour),
			},
			stats: entity.ContentStats{
				Views: 0,
				Likes: 0,
			},
			expected: entity.ScoreComponents{
				BaseScore:       0.0,
				TypeMultiplier:  1.5,
				RecencyScore:    10.0,
				EngagementScore: 0.0,
				// Final: (0 * 1.5) + 10 + 0 = 10
				FinalScore: 10.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.Calculate(tt.content, tt.stats)
			
			// Use epsilon for float comparison
			epsilon := 0.0001
			assert.InDelta(t, tt.expected.BaseScore, result.BaseScore, epsilon, "BaseScore mismatch")
			assert.InDelta(t, tt.expected.TypeMultiplier, result.TypeMultiplier, epsilon, "TypeMultiplier mismatch")
			assert.InDelta(t, tt.expected.RecencyScore, result.RecencyScore, epsilon, "RecencyScore mismatch")
			assert.InDelta(t, tt.expected.EngagementScore, result.EngagementScore, epsilon, "EngagementScore mismatch")
			assert.InDelta(t, tt.expected.FinalScore, result.FinalScore, epsilon, "FinalScore mismatch")
		})
	}
}

func TestScoringService_Recency(t *testing.T) {
	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	timeProvider := func() time.Time { return now }
	
	config := entity.ScoringConfig{
		RecencyWeekScore:    10.0,
		RecencyMonthScore:   5.0,
		RecencyQuarterScore: 2.0,
	}
	
	service := NewScoringService(config, timeProvider)
	
	tests := []struct {
		name        string
		publishedAt time.Time
		expected    float64
	}{
		{"Today", now, 10.0},
		{"6 Days Ago", now.Add(-6 * 24 * time.Hour), 10.0},
		{"8 Days Ago", now.Add(-8 * 24 * time.Hour), 5.0},
		{"29 Days Ago", now.Add(-29 * 24 * time.Hour), 5.0},
		{"31 Days Ago", now.Add(-31 * 24 * time.Hour), 2.0},
		{"89 Days Ago", now.Add(-89 * 24 * time.Hour), 2.0},
		{"91 Days Ago", now.Add(-91 * 24 * time.Hour), 0.0},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := entity.Content{PublishedAt: tt.publishedAt}
			// Access private method via Calculate or just trust Calculate uses it?
			// Since we are testing public API, we check the component in Calculate result.
			// But to isolate, we can rely on Calculate's RecencyScore field.
			
			// We need dummy stats/type to call Calculate
			res := service.Calculate(content, entity.ContentStats{})
			assert.Equal(t, tt.expected, res.RecencyScore)
		})
	}
}
