package config

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
)

// DatabaseConfigProvider wraps the original Viper config but overrides GetScoringConfig
// to fetch from the database.
type DatabaseConfigProvider struct {
	baseProvider ports.ConfigProvider
	repo         ports.ScoringRepository // We need a new port for this
}

func NewDatabaseConfigProvider(base ports.ConfigProvider, repo ports.ScoringRepository) *DatabaseConfigProvider {
	return &DatabaseConfigProvider{
		baseProvider: base,
		repo:         repo,
	}
}

// Delegate basic config to Viper
func (p *DatabaseConfigProvider) GetAppConfig() *entity.AppConfig {
	return p.baseProvider.GetAppConfig()
}

// Override Scoring Config to read from DB
func (p *DatabaseConfigProvider) GetScoringConfig() entity.ScoringConfig {
	// Default config as fallback
	config := entity.ScoringConfig{
		VideoTypeMultiplier:    1.5,
		TextTypeMultiplier:     1.0,
		RecencyWeekScore:       5.0,
		RecencyMonthScore:      3.0,
		RecencyQuarterScore:    1.0,
		VideoEngagementWeight:  10.0,
		TextEngagementWeight:   5.0,
		VideoViewsDivisor:      1000.0,
		VideoLikesDivisor:      100.0,
		TextReadingTimeDivisor: 1.0,
		TextReactionsDivisor:   50.0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rules, err := p.repo.GetScoringRules(ctx)
	if err != nil {
		return config
	}

	type VideoConfig struct {
		TypeMultiplier   float64 `json:"type_multiplier"`
		EngagementWeight float64 `json:"engagement_weight"`
		ViewsDivisor     float64 `json:"views_divisor"`
		LikesDivisor     float64 `json:"likes_divisor"`
	}
	type ArticleConfig struct {
		TypeMultiplier     float64 `json:"type_multiplier"`
		EngagementWeight   float64 `json:"engagement_weight"`
		ReadingTimeDivisor float64 `json:"reading_time_divisor"`
		ReactionsDivisor   float64 `json:"reactions_divisor"`
	}
	type RecencyConfig struct {
		WeekScore    float64 `json:"week_score"`
		MonthScore   float64 `json:"month_score"`
		QuarterScore float64 `json:"quarter_score"`
	}

	for key, value := range rules {
		switch key {
		case "video_config":
			var vc VideoConfig
			if err := json.Unmarshal(value, &vc); err == nil {
				config.VideoTypeMultiplier = vc.TypeMultiplier
				config.VideoEngagementWeight = vc.EngagementWeight
				config.VideoViewsDivisor = vc.ViewsDivisor
				config.VideoLikesDivisor = vc.LikesDivisor
			}
		case "article_config":
			var ac ArticleConfig
			if err := json.Unmarshal(value, &ac); err == nil {
				config.TextTypeMultiplier = ac.TypeMultiplier
				config.TextEngagementWeight = ac.EngagementWeight
				config.TextReadingTimeDivisor = ac.ReadingTimeDivisor
				config.TextReactionsDivisor = ac.ReactionsDivisor
			}
		case "recency_config":
			var rc RecencyConfig
			if err := json.Unmarshal(value, &rc); err == nil {
				config.RecencyWeekScore = rc.WeekScore
				config.RecencyMonthScore = rc.MonthScore
				config.RecencyQuarterScore = rc.QuarterScore
			}
		}
	}

	return config
}
