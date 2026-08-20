package api

import (
	"fmt"
	"net/http"

	"github.com/TNJKL/bookmark-management/docs"
	_ "github.com/TNJKL/bookmark-management/docs" // Load tài liệu Swagger đã generate
	"github.com/TNJKL/bookmark-management/internal/api/middleware"
	"github.com/TNJKL/bookmark-management/internal/app/handler/bookmark"
	genpassHandler "github.com/TNJKL/bookmark-management/internal/app/handler/genpass"
	"github.com/TNJKL/bookmark-management/internal/app/handler/healthcheck"
	"github.com/TNJKL/bookmark-management/internal/app/handler/link"
	userHandler "github.com/TNJKL/bookmark-management/internal/app/handler/user"
	bookmarkRepo "github.com/TNJKL/bookmark-management/internal/app/repository/bookmark"
	"github.com/TNJKL/bookmark-management/internal/app/repository/cache"
	"github.com/TNJKL/bookmark-management/internal/app/repository/ping"
	"github.com/TNJKL/bookmark-management/internal/app/repository/queue"
	"github.com/TNJKL/bookmark-management/internal/app/repository/ratelimit"
	"github.com/TNJKL/bookmark-management/internal/app/repository/urlstorage"
	"github.com/TNJKL/bookmark-management/internal/app/repository/user"
	bookmarkSvc "github.com/TNJKL/bookmark-management/internal/app/service/bookmark"
	"github.com/TNJKL/bookmark-management/internal/app/service/genpass"
	healthcheck2 "github.com/TNJKL/bookmark-management/internal/app/service/healthcheck"
	linkSvc "github.com/TNJKL/bookmark-management/internal/app/service/link"
	queueSvc "github.com/TNJKL/bookmark-management/internal/app/service/queue"
	userSvc "github.com/TNJKL/bookmark-management/internal/app/service/user"
	"github.com/TNJKL/bookmark-management/pkg/jwtutils"
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

const queueName = "bookmark-import"

// Struct thực tế implement interface
type engine struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
	db          *gorm.DB
	jwtGen      jwtutils.JWTGenerator
	jwtVal      jwtutils.JWTValidator
}

// EngineOpts holds the configuration and dependencies required to initialize the API Engine
type EngineOpts struct {
	App         *gin.Engine
	Cfg         *Config
	RedisClient *redis.Client
	Db          *gorm.DB
	JWTGen      jwtutils.JWTGenerator
	JWTVal      jwtutils.JWTValidator
}

// NewEngine creates and configures a new HTTP API Engine.
func NewEngine(opts *EngineOpts) Engine {

	app := &engine{
		app:         opts.App, // Tạo Gin router
		cfg:         opts.Cfg,
		redisClient: opts.RedisClient,
		db:          opts.Db,
		jwtGen:      opts.JWTGen,
		jwtVal:      opts.JWTVal,
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
	genPassHandler     genpassHandler.GenPass
	healthCheckHandler healthcheck.HealthCheck
	urlStorageHandler  link.ShortenURL
	userHandler        userHandler.Handler
	bookmarkHandler    bookmark.Handler
}

// initHandlers initializes the api handlers
func (e *engine) initHandlers() *handlers {
	genPassSvc := genpass.NewGenPass()
	keyGen := utils.NewKeyGenerator()
	pingRepo := ping.NewHealthRepository(e.redisClient)
	healthCheckSvc := healthcheck2.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceID, pingRepo)
	urlStorage := urlstorage.NewURLStorage(e.redisClient)
	bookmarkRepository := bookmarkRepo.NewRepository(e.db)
	shortenUrlSvc := linkSvc.NewShortenUrl(urlStorage, keyGen, bookmarkRepository)

	//init user handler
	userRepo := user.NewSQLRepository(e.db)
	hasher := utils.NewHasher()
	userService := userSvc.NewService(userRepo, hasher, e.jwtGen)

	//cache repo
	cacheRepo := cache.NewRedisDB(e.redisClient)

	//init Base62
	base62 := utils.NewBase62(e.cfg.Base62XORSecret)

	//init queue repo and service
	queueRepo := queue.NewRedisQueueRepo(e.redisClient, queueName)
	queueService := queueSvc.NewService(queueRepo)

	//init bookmark
	bookmarkService := bookmarkSvc.NewService(bookmarkRepository, base62, e.db)
	bookmarkServiceWithCache := bookmarkSvc.NewServiceWithCache(bookmarkService, cacheRepo)

	return &handlers{
		genPassHandler:     genpassHandler.NewGenPass(genPassSvc),
		healthCheckHandler: healthcheck.NewHealthCheck(healthCheckSvc),
		urlStorageHandler:  link.NewShortenURL(shortenUrlSvc),
		userHandler:        userHandler.NewHandler(userService),
		bookmarkHandler:    bookmark.NewHandler(bookmarkServiceWithCache, queueService),
	}
}

// middlewares represents the struct for containing all the necessary middlewares for API
type middlewares struct {
	jwtAuth   middleware.JWTAuth
	ratelimit middleware.RateLimit
}

// initMiddlewares initializes the api middlewares
func (e *engine) initMiddlewares() middlewares {
	rateLimitRepo := ratelimit.NewRedisRepo(e.redisClient)
	return middlewares{
		jwtAuth:   middleware.NewJWTAuth(e.jwtVal),
		ratelimit: middleware.NewRateLimit(rateLimitRepo),
	}
}

// initRoutes initializes the api routes
func (e *engine) initRoutes() {
	allHandler := e.initHandlers()
	allMiddlewares := e.initMiddlewares()

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
		v1Routes.POST("/users/login", allHandler.userHandler.Login)

	}
	//private routes (need Auth)
	privateRoutes := e.app.Group("")
	privateRoutes.Use(allMiddlewares.jwtAuth.JWTAuth())
	privateRoutes.Use(allMiddlewares.ratelimit.RateLimit())
	{
		privateV1Routes := privateRoutes.Group("/v1")
		{
			//self endpoints
			privateV1Routes.GET("self/info", allHandler.userHandler.GetSelfInfo)
			privateV1Routes.PUT("self/info", allHandler.userHandler.UpdateSelfInfo)

			//bookmark endpoints
			privateV1Routes.POST("/bookmarks", allHandler.bookmarkHandler.CreateBookmark)
			privateV1Routes.GET("/bookmarks", allHandler.bookmarkHandler.GetBookmarks)
			privateV1Routes.PUT("/bookmarks/:id", allHandler.bookmarkHandler.UpdateBookmark)
			privateV1Routes.DELETE("/bookmarks/:id", allHandler.bookmarkHandler.DeleteBookmark)
			privateV1Routes.POST("/bookmarks/import", allHandler.bookmarkHandler.ImportBookmarks)

		}

	}
}
