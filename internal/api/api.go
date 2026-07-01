package api

import (
	"fmt"
	"net/http"

	"github.com/TNJKL/bookmark-management/internal/handler"
	"github.com/TNJKL/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
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

}
