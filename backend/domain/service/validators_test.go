package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagNormalizer_Normalize(t *testing.T) {
	normalizer := NewTagNormalizer()

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Mixed Case and Whitespace",
			input:    []string{" Go ", "python", "JAVA "},
			expected: []string{"go", "python", "java"},
		},
		{
			name:     "Duplicates",
			input:    []string{"go", "Go", "GO"},
			expected: []string{"go"},
		},
		{
			name:     "Empty Strings",
			input:    []string{"", "  ", "valid"},
			expected: []string{"valid"},
		},
		{
			name:     "Empty Input",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizer.Normalize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchQueryValidator_ValidatePageSize(t *testing.T) {
	validator := NewSearchQueryValidator(50)

	tests := []struct {
		name     string
		pageSize int32
		wantErr  bool
	}{
		{"Valid", 10, false},
		{"Max", 50, false},
		{"Zero", 0, true},
		{"Negative", -1, true},
		{"Exceeds Max", 51, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePageSize(tt.pageSize)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearchQueryValidator_ValidatePage(t *testing.T) {
	validator := NewSearchQueryValidator(50)

	tests := []struct {
		name    string
		page    int32
		wantErr bool
	}{
		{"Valid", 1, false},
		{"Zero", 0, true},
		{"Negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePage(tt.page)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
