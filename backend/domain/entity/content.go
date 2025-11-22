package entity

import "time"

type ContentType string

const (
	ContentTypeVideo   ContentType = "video"
	ContentTypeArticle ContentType = "article"
)

type Content struct {
	ID                int64
	ProviderID        int64
	ProviderContentID string
	Title             string
	ContentType       ContentType
	PublishedAt       time.Time
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (c Content) IsVideo() bool {
	return c.ContentType == ContentTypeVideo
}

func (c Content) IsArticle() bool {
	return c.ContentType == ContentTypeArticle
}
