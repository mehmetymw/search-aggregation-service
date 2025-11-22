package ports

import "github.com/mehmetymw/search-aggregation-service/backend/domain/entity"

type ConfigProvider interface {
	GetAppConfig() *entity.AppConfig
}
