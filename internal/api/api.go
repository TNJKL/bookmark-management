package api

import (
	"fmt"
	"net/http"

	_ "github.com/TNJKL/bookmark-management/docs" // Load tài liệu Swagger đã generate
	"github.com/TNJKL/bookmark-management/internal/handler"
	"github.com/TNJKL/bookmark-management/internal/repository"
	"github.com/TNJKL/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
}

// NewEngine creates and configures a new HTTP API Engine.
func NewEngine(cfg *Config, redis *redis.Client) Engine {

	app := &engine{
		app:         gin.Default(), // Tạo Gin router
		cfg:         cfg,
		redisClient: redis,
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

func (e *engine) initRoutes() {
	//khai bao genpass handler

	// Bước 1: Tạo Service (tầng logic)
	genPassSvc := service.NewGenPass()

	// Bước 2: Tạo Handler, TRUYỀN service vào (DI)
	genPassHandler := handler.NewGenPass(genPassSvc)

	//khai bao Health check handler

	// Bước 1: Tạo Service
	healthRepo := repository.NewHealthRepository(e.redisClient)
	healthCheckSvc := service.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceID, healthRepo)

	// Bước 2: Tạo Handler, TRUYỀN service vào (DI)
	healthCheckHandler := handler.NewHealthCheck(healthCheckSvc)

	//khai bao Shorten URL handler

	//tạo repository
	urlStorage := repository.NewURLStorage(e.redisClient)

	// Bước 1: Tạo Service
	shortenUrlSvc := service.NewShortenUrl(urlStorage, genPassSvc)

	// Bước 2: Tạo Handler, TRUYỀN service vào (DI)
	urlStorageHandler := handler.NewShortenURL(shortenUrlSvc)

	e.app.GET("/genpass", genPassHandler.GeneratePassword)
	e.app.GET("/health-check", healthCheckHandler.HealthCheck)
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.app.POST("/v1/links/shorten", urlStorageHandler.ShortenLink)

}
