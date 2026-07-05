package main

import (
	"context"
	"time"

	"github.com/TNJKL/bookmark-management/pkg/redis"
)

func main() {
	rclient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}
	rclient.Set(context.Background(), "GALA", "facebook.com", time.Hour)

	rclient2, err := redis.NewClient("CACHE")
	if err != nil {
		panic(err)
	}

	rclient2.Set(context.Background(), "Kakarot", "amazom.com", time.Hour)
	rclient2.Set(context.Background(), "Vegeta", "youtube.com", time.Hour)
	rclient2.Set(context.Background(), "Gohan", "instagram.com", time.Hour)

}
