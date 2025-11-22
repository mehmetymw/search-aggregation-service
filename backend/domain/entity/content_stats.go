package entity

import "time"

type ContentStats struct {
	ContentID   int64
	Views       int64
	Likes       int64
	DurationSec int32
	ReadingTime int32
	Reactions   int64
	Comments    int64
	LastSyncAt  time.Time
}
