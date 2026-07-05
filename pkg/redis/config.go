package redis

import "github.com/kelseyhightower/envconfig"

// Không viết hoa ( để tất cả private ) chỉ dùng trong package "redis" , ko muốn người khác biết các thông tin nhạy cảm
// Chỉ cần gọi newConfig ra là được
type config struct {
	Address  string `default:"localhost:6379" envconfig:"REDIS_ADDRESS"`
	Password string `default:"" envconfig:"REDIS_PASSWORD"`
	DB       int    `default:"0" envconfig:"REDIS_DB"`
}

// newConfig create a new config struct based on the environment variables
func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
