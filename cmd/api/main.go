package main

import (
	"github.com/TNJKL/bookmark-management/internal/api"
	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
)

// @title       Bookmark Management API
// @version     1.0.0
// @description API Swagger for Bookmark-Management.
// @BasePath    /
func main() {
	//create app config
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	//create redis client
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}
	app := api.NewEngine(cfg, redisClient) //Buoc 1 : khoi tao Engine (Khoi tao moi thu)
	err = app.Start()                      //Buoc 2 :Chay server
	if err != nil {
		panic(err)
	}
}
