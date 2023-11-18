package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-go:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		visits, err := rdb.Get(ctx, "visits").Int()
		if err != nil {
			visits = 0
		}
		fmt.Fprintf(w, "Number of visits is %d", visits)
		rdb.Set(ctx, "visits", visits+1, 0)
	})

	rdb.Set(ctx, "visits", 0, 0)
	http.ListenAndServe(":8080", nil)
}
