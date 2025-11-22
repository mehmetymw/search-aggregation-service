package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
	db "github.com/mehmetymw/search-aggregation-service/backend/db/generated"
)

type ScoringRepository struct {
	queries *db.Queries
}

func NewScoringRepository(conn *sql.DB) ports.ScoringRepository {
	return &ScoringRepository{
		queries: db.New(conn),
	}
}

func (r *ScoringRepository) GetScoringRules(ctx context.Context) (map[string][]byte, error) {
	rules, err := r.queries.GetScoringRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("get scoring rules: %w", err)
	}

	result := make(map[string][]byte)
	for _, rule := range rules {
		result[rule.Key] = rule.Value
	}
	return result, nil
}
