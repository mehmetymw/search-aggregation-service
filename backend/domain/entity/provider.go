package entity

import "time"

const (
	ProviderFormatJSON = "json"
	ProviderFormatXML  = "xml"
)

type Provider struct {
	ID        int64
	Name      string
	Code      string
	Format    string
	BaseURL   string
	IsEnabled bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
