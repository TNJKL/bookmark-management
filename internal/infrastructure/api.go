package infrastructure

import (
	"github.com/TNJKL/bookmark-management/internal/api"
	"github.com/TNJKL/bookmark-management/pkg/logger"
	"github.com/TNJKL/bookmark-management/pkg/utils"
	"github.com/gin-gonic/gin"
)

// CreateAPIConfig creates API configuration based on environment variables
func CreateAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	utils.NoErr(err)
	return cfg
}

// CreateAPI creates API application with default setup
func CreateAPI() api.Engine {
	// Init config
	cfg := CreateAPIConfig()

	//Init logger
	logger.SetLogLevel(cfg.LogLevel)

	// Init redis
	redisClient := CreateRedisConn()

	// Init sql db
	sqlDB := CreateSQLDBWithMigration()

	// Init jwt gen and validator
	jwtGen, jwtVal := CreateJWTProvider()

	app := gin.New()

	return api.NewEngine(&api.EngineOpts{
		App:         app,
		Cfg:         cfg,
		RedisClient: redisClient,
		Db:          sqlDB,
		JWTGen:      jwtGen,
		JWTVal:      jwtVal,
	})
}
