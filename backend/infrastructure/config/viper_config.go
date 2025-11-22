package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
)

type ViperConfig struct {
	v *viper.Viper
}

func NewViperConfig(configPath string) (ports.ConfigProvider, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return &ViperConfig{v: v}, nil
}

func (c *ViperConfig) GetAppConfig() *entity.AppConfig {
	var config entity.AppConfig
	if err := c.v.Unmarshal(&config); err != nil {
		return &config
	}
	return &config
}