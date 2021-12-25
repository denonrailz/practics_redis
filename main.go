package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

//Basic usage of redis, just create connections setting and getting sample data
func main() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	pass := os.Getenv("REDIS_PASS")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", host, port),
		Password: pass, // no password set
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DB: 0, // use default DB
	})

	err := rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(context.Background(), "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}
