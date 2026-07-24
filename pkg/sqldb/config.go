package sqldb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// config stores the database connection settings read from environment variables
type config struct {
	Host     string `default:"localhost" envconfig:"DB_HOST"`
	User     string `default:"admin" envconfig:"DB_USER"`
	Password string `default:"admin" envconfig:"DB_PASSWORD"`
	DBName   string `default:"bookmark" envconfig:"DB_NAME"`
	Port     string `default:"5432" envconfig:"DB_PORT"`
	SSLMode  string `default:"disable" envconfig:"DB_SSL_MODE"`
	Timezone string `default:"UTC" envconfig:"DB_TIMEZONE"`
}

// newConfig creates a new config struct based on the environment variables
func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// GetDSN returns the Data Source Name (DSN) string for the PostgreSQL database connection
func (cfg *config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode, cfg.Timezone)
}
