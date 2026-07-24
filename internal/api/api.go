package api

import (
	"fmt"
	"net/http"

	"github.com/TNJKL/bookmark-management/docs"
	_ "github.com/TNJKL/bookmark-management/docs" // Load tài liệu Swagger đã generate
	"github.com/TNJKL/bookmark-management/internal/handler"
	userHandler "github.com/TNJKL/bookmark-management/internal/handler/user"
	"github.com/TNJKL/bookmark-management/internal/repository/ping"
	"github.com/TNJKL/bookmark-management/internal/repository/urlstorage"
	"github.com/TNJKL/bookmark-management/internal/repository/user"
	"github.com/TNJKL/bookmark-management/internal/service"
	userSvc "github.com/TNJKL/bookmark-management/internal/service/user"
	"github.com/TNJKL/bookmark-management/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Interface định nghĩa "Engine có thể làm gì"
// Engine defines the contract for starting the server and handling HTTP requests.
type Engine interface {
	Start() error
	ServerHTTP(w http.ResponseWriter, req *http.Request)
}

// Struct thực tế implement interface
type engine struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
	db          *gorm.DB
}

// NewEngine creates and configures a new HTTP API Engine.
func NewEngine(cfg *Config, redis *redis.Client, db *gorm.DB) Engine {

	app := &engine{
		app:         gin.Default(), // Tạo Gin router
		cfg:         cfg,
		redisClient: redis,
		db:          db,
	}
	app.initRoutes() // Đăng ký các routes

	return app
}

// Start runs the HTTP server on the configured port.
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.Apport))
}

// Server HTTP to test the API endpoint
// ServerHTTP handles HTTP requests directly, primarily used for testing endpoints.
func (e *engine) ServerHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

// handlers aggregates all HTTP handler dependencies used to register
// the application's routes.
type handlers struct {
	genPassHandler     handler.GenPass
	healthCheckHandler handler.HealthCheck
	urlStorageHandler  handler.ShortenURL
	userHandler        userHandler.Handler
}

// initHandlers initializes the api handlers
func (e *engine) initHandlers() *handlers {
	genPassSvc := service.NewGenPass()
	keyGen := utils.NewKeyGenerator()
	pingRepo := ping.NewHealthRepository(e.redisClient)
	healthCheckSvc := service.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceID, pingRepo)
	urlStorage := urlstorage.NewURLStorage(e.redisClient)
	shortenUrlSvc := service.NewShortenUrl(urlStorage, keyGen)

	//user
	userRepo := user.NewSQLRepository(e.db)
	hasher := utils.NewHasher()
	userSvc := userSvc.NewService(userRepo, hasher)

	return &handlers{
		genPassHandler:     handler.NewGenPass(genPassSvc),
		healthCheckHandler: handler.NewHealthCheck(healthCheckSvc),
		urlStorageHandler:  handler.NewShortenURL(shortenUrlSvc),
		userHandler:        userHandler.NewHandler(userSvc),
	}
}

// initRoutes initializes the api routes
func (e *engine) initRoutes() {
	allHandler := e.initHandlers()

	//genpass
	e.app.GET("/genpass", allHandler.genPassHandler.GeneratePassword)

	//health-check
	e.app.GET("/health-check", allHandler.healthCheckHandler.HealthCheck)

	//Init swagger routes
	docs.SwaggerInfo.BasePath = e.cfg.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1Routes := e.app.Group("/v1")
	{
		//link-related
		v1Routes.POST("/links/shorten", allHandler.urlStorageHandler.ShortenLink)
		v1Routes.GET("/links/redirect/:code", allHandler.urlStorageHandler.Redirect)

		//user
		v1Routes.POST("/users/register", allHandler.userHandler.Register)
	}

}
