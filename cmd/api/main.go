package main

import (
	"log"

	"github.com/TNJKL/bookmark-management/internal/infrastructure"
)

// @title       Bookmark Management API
// @version     4.0.0
// @description API Swagger for Bookmark-Management.
// @BasePath    /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	//Init api
	a := infrastructure.CreateAPI()

	//Start api app
	if err := a.Start(); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
