package service

import (
	"fmt"
	"strings"
)

type TagNormalizer struct{}

func NewTagNormalizer() *TagNormalizer {
	return &TagNormalizer{}
}

func (tn *TagNormalizer) Normalize(tags []string) []string {
	normalized := make([]string, 0, len(tags))
	seen := make(map[string]bool)
	
	for _, tag := range tags {
		cleaned := strings.TrimSpace(strings.ToLower(tag))
		if cleaned != "" && !seen[cleaned] {
			normalized = append(normalized, cleaned)
			seen[cleaned] = true
		}
	}
	
	return normalized
}

type SearchQueryValidator struct {
	maxPageSize int
}

func NewSearchQueryValidator(maxPageSize int) *SearchQueryValidator {
	return &SearchQueryValidator{
		maxPageSize: maxPageSize,
	}
}

func (v *SearchQueryValidator) ValidatePageSize(pageSize int32) error {
	if pageSize <= 0 {
		return fmt.Errorf("page size must be greater than 0")
	}
	if pageSize > int32(v.maxPageSize) {
		return fmt.Errorf("page size exceeds maximum of %d", v.maxPageSize)
	}
	return nil
}

func (v *SearchQueryValidator) ValidatePage(page int32) error {
	if page <= 0 {
		return fmt.Errorf("page must be greater than 0")
	}
	return nil
}
