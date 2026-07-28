package main

import (
	"github.com/TNJKL/bookmark-management/internal/api"
	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/pkg/jwtutils"
	"github.com/TNJKL/bookmark-management/pkg/logger"
	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
	"github.com/TNJKL/bookmark-management/pkg/sqldb"
	"github.com/gin-gonic/gin"
)

// @title       Bookmark Management API
// @version     4.0.0
// @description API Swagger for Bookmark-Management.
// @BasePath    /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	//create app config
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	//set log level
	logger.SetLogLevel(cfg.LogLevel)

	//init DB
	db, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}
	//auto migrate
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	//create redis client
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}

	//create jwtGen
	jwtGen, err := jwtutils.NewJWTGenerator("./private.pem")
	if err != nil {
		panic(err)
	}

	//create jwtVal
	jwtVal, err := jwtutils.NewJWTValidator("./public.pem")
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(&api.EngineOpts{
		App:         gin.New(),
		Cfg:         cfg,
		RedisClient: redisClient,
		Db:          db,
		JWTGen:      jwtGen,
		JWTVal:      jwtVal,
	}) //Buoc 1 : khoi tao Engine (Khoi tao moi thu)
	err = app.Start() //Buoc 2 :Chay server
	if err != nil {
		panic(err)
	}
}
