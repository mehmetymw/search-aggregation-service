package service

import (
	"time"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
)

type TimeProvider func() time.Time

type ScoringService struct {
	config       entity.ScoringConfig
	timeProvider TimeProvider
}

func NewScoringService(config entity.ScoringConfig, timeProvider TimeProvider) *ScoringService {
	return &ScoringService{
		config:       config,
		timeProvider: timeProvider,
	}
}

func (s *ScoringService) Calculate(content entity.Content, stats entity.ContentStats) entity.ScoreComponents {
	now := s.timeProvider()
	
	baseScore := s.computeBaseScore(content, stats)
	typeMultiplier := s.getTypeMultiplier(content)
	recencyScore := s.computeRecencyScore(content, now)
	engagementScore := s.computeEngagementScore(content, stats)
	
	finalScore := (baseScore * typeMultiplier) + recencyScore + engagementScore
	
	return entity.ScoreComponents{
		BaseScore:       baseScore,
		TypeMultiplier:  typeMultiplier,
		RecencyScore:    recencyScore,
		EngagementScore: engagementScore,
		FinalScore:      finalScore,
	}
}

func (s *ScoringService) computeBaseScore(content entity.Content, stats entity.ContentStats) float64 {
	switch content.ContentType {
	case entity.ContentTypeVideo:
		viewsDivisor := s.config.VideoViewsDivisor
		if viewsDivisor == 0 {
			viewsDivisor = 1.0
		}
		likesDivisor := s.config.VideoLikesDivisor
		if likesDivisor == 0 {
			likesDivisor = 1.0
		}
		
		viewsScore := float64(stats.Views) / viewsDivisor
		likesScore := float64(stats.Likes) / likesDivisor
		return viewsScore + likesScore
		
	case entity.ContentTypeArticle:
		readingTimeDivisor := s.config.TextReadingTimeDivisor
		if readingTimeDivisor == 0 {
			readingTimeDivisor = 1.0
		}
		reactionsDivisor := s.config.TextReactionsDivisor
		if reactionsDivisor == 0 {
			reactionsDivisor = 1.0
		}
		
		readingTimeScore := float64(stats.ReadingTime) / readingTimeDivisor
		reactionsScore := float64(stats.Reactions) / reactionsDivisor
		return readingTimeScore + reactionsScore
		
	default:
		return 0.0
	}
}

func (s *ScoringService) getTypeMultiplier(content entity.Content) float64 {
	switch content.ContentType {
	case entity.ContentTypeVideo:
		return s.config.VideoTypeMultiplier
	case entity.ContentTypeArticle:
		return s.config.TextTypeMultiplier
	default:
		return 1.0
	}
}

func (s *ScoringService) computeRecencyScore(content entity.Content, now time.Time) float64 {
	elapsed := now.Sub(content.PublishedAt)
	daysSincePublish := elapsed.Hours() / 24.0
	
	const (
		week    = 7.0
		month   = 30.0
		quarter = 90.0
	)
	
	if daysSincePublish <= week {
		return s.config.RecencyWeekScore
	} else if daysSincePublish <= month {
		return s.config.RecencyMonthScore
	} else if daysSincePublish <= quarter {
		return s.config.RecencyQuarterScore
	}
	
	return 0.0
}

func (s *ScoringService) computeEngagementScore(content entity.Content, stats entity.ContentStats) float64 {
	switch content.ContentType {
	case entity.ContentTypeVideo:
		if stats.Views == 0 {
			return 0.0
		}
		ratio := float64(stats.Likes) / float64(stats.Views)
		return ratio * s.config.VideoEngagementWeight
		
	case entity.ContentTypeArticle:
		if stats.ReadingTime == 0 {
			return 0.0
		}
		ratio := float64(stats.Reactions) / float64(stats.ReadingTime)
		return ratio * s.config.TextEngagementWeight
		
	default:
		return 0.0
	}
}
