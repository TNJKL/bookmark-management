package main

import (
	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/pkg/sqldb"
)

func main() {
	dbClient, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}

	dbClient.AutoMigrate(&model.User{})
	dbClient.Create(&model.User{
		Username: "alice",
		Password: "123456",
	})
}
