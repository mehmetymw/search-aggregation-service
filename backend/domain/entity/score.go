package entity

type ScoreComponents struct {
	BaseScore        float64
	TypeMultiplier   float64
	RecencyScore     float64
	EngagementScore  float64
	FinalScore       float64
}

type ScoringConfig struct {
	VideoTypeMultiplier    float64
	TextTypeMultiplier     float64
	RecencyWeekScore       float64
	RecencyMonthScore      float64
	RecencyQuarterScore    float64
	VideoEngagementWeight  float64
	TextEngagementWeight   float64
	VideoViewsDivisor      float64
	VideoLikesDivisor      float64
	TextReadingTimeDivisor float64
	TextReactionsDivisor   float64
}
