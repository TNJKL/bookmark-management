package main

import "github.com/TNJKL/bookmark-management/internal/api"

func main() {
	app := api.NewEngine()
	err := app.Start()
	if err != nil {
		panic(err)
	}
}
