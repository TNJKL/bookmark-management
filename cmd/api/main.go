package main

import "github.com/TNJKL/bookmark-management/internal/api"

// @title           Bookmark Management API
// @version         1.0
// @description     API Swagger của Bookmark-Management.
// @host            localhost:8080
// @BasePath        /
func main() {
	//create app config
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg) //Buoc 1 : khoi tao Engine (Khoi tao moi thu)
	err = app.Start()         //Buoc 2 :Chay server
	if err != nil {
		panic(err)
	}
}
