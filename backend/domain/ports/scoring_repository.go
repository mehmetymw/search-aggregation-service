package ports

import "context"

type ScoringRepository interface {
	GetScoringRules(ctx context.Context) (map[string][]byte, error)
}
