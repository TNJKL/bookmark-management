package api

import (
	"fmt"
	"net/http"

	_ "github.com/TNJKL/bookmark-management/docs" // Load tài liệu Swagger đã generate
	"github.com/TNJKL/bookmark-management/internal/handler"
	"github.com/TNJKL/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Interface định nghĩa "Engine có thể làm gì"
type Engine interface {
	Start() error
	ServerHTTP(w http.ResponseWriter, req *http.Request)
}

// Struct thực tế implement interface
type engine struct {
	app *gin.Engine
	cfg *Config
}

func NewEngine(cfg *Config) Engine {

	app := &engine{
		app: gin.Default(), // Tạo Gin router
		cfg: cfg,
	}
	app.initRoutes() // Đăng ký các routes

	return app
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.Apport))
}

// Server HTTP to test the API endpoint
func (e *engine) ServerHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

func (e *engine) initRoutes() {
	//khai bao genpass handler

	// Bước 1: Tạo Service (tầng logic)
	genPassSvc := service.NewGenPass()

	// Bước 2: Tạo Handler, TRUYỀN service vào (DI)
	genPassHandler := handler.NewGenPass(genPassSvc)

	// Bước 3: Gắn Handler vào route
	e.app.GET("/genpass", genPassHandler.GeneratePassword)

	//khai bao Health check handler

	// Bước 1: Tạo Service
	healthCheckSvc := service.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceID)

	// Bước 2: Tạo Handler, TRUYỀN service vào (DI)
	healthCheckHandler := handler.NewHealthCheck(healthCheckSvc)

	// Bước 3: Gắn Handler vào route
	e.app.GET("/health-check", healthCheckHandler.HealthCheck)

	//phần này em nhờ AI Gen
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
