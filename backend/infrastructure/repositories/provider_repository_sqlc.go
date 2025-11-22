package repositories

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/mehmetymw/search-aggregation-service/backend/db/generated"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
)

type ProviderRepositorySqlc struct {
	db      *sql.DB
	queries *db.Queries
}

func NewProviderRepository(database *sql.DB) ports.ProviderRepository {
	return &ProviderRepositorySqlc{
		db:      database,
		queries: db.New(database),
	}
}

func (r *ProviderRepositorySqlc) GetAllEnabled(ctx context.Context) ([]entity.Provider, error) {
	rows, err := r.queries.GetAllEnabledProviders(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all enabled providers: %w", err)
	}

	providers := make([]entity.Provider, 0, len(rows))
	for _, row := range rows {
		providers = append(providers, dbRowToProvider(row))
	}

	return providers, nil
}

func (r *ProviderRepositorySqlc) GetByCode(ctx context.Context, code string) (*entity.Provider, error) {
	row, err := r.queries.GetProviderByCode(ctx, code)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get provider by code: %w", err)
	}

	provider := dbRowToProvider(row)
	return &provider, nil
}

func (r *ProviderRepositorySqlc) GetByID(ctx context.Context, id int64) (*entity.Provider, error) {
	row, err := r.queries.GetProviderByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get provider by id: %w", err)
	}

	provider := dbRowToProvider(row)
	return &provider, nil
}

func (r *ProviderRepositorySqlc) UpsertProvider(ctx context.Context, provider entity.Provider) error {
	err := r.queries.UpsertProvider(ctx, db.UpsertProviderParams{
		Name:      provider.Name,
		Code:      provider.Code,
		Format:    provider.Format,
		BaseUrl:   provider.BaseURL,
		IsEnabled: provider.IsEnabled,
	})
	if err != nil {
		return fmt.Errorf("upsert provider: %w", err)
	}
	return nil
}

func dbRowToProvider(row db.Provider) entity.Provider {
	var formatStr string
	switch v := row.Format.(type) {
	case string:
		formatStr = v
	case []byte:
		formatStr = string(v)
	default:
		formatStr = fmt.Sprintf("%v", row.Format)
	}

	return entity.Provider{
		ID:        row.ID,
		Name:      row.Name,
		Code:      row.Code,
		Format:    formatStr,
		BaseURL:   row.BaseUrl,
		IsEnabled: row.IsEnabled,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
