package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

func getRedisClient() (*redis.Client, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	pass := os.Getenv("REDIS_PASS")

	cl := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", host, port),
		Password: pass, // no password set
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DB: 0, // use default DB
	})
	_, err := cl.Ping(context.Background()).Result()
	return cl, err
}
