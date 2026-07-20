package api

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Apport      string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark_service'" envconfig:"SERVICE_NAME"`
	InstanceID  string `default:"" envconfig:"INSTANCE_ID"`
	LogLevel    string `default:"info" envconfig:"LOG_LEVEL"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("api", cfg)
	if err != nil {
		return nil, err
	}

	//Kiểm tra nếu instanceId "" thì tạo random ra 1 UUID
	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}
	return cfg, err
}
