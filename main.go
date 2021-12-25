package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type rateLimiter struct {
	ttl       time.Duration //
	threshold int           // rate limiting threshold
	client    *redis.Client
}

func (rl *rateLimiter) Rate(ctx context.Context, key string) error {
	val, _ := rl.client.Get(ctx, key).Int()
	if val > rl.threshold {
		return errors.New("max rate limiter reached, please try againg later")
	}
	rl.client.Incr(ctx, key)
	rl.client.Expire(ctx, key, rl.ttl)
	return nil
}

func main() {
	redisClient, err := getRedisClient()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to redis: success")
	rLimiter := rateLimiter{
		ttl:       time.Second * 10,
		threshold: 10,
		client:    redisClient,
	}
	ctx := context.Background()

	var i, j int
	for {
		if i > 50 {
			break
		}
		err := rLimiter.Rate(ctx, "sample")
		if err != nil {
			fmt.Printf("%v. Try: %v. Wait...\n", err, j)
			time.Sleep(100 * time.Millisecond)
			j++
		} else {
			fmt.Printf("%v. Some action\n", i)
			j = 0
			i++
		}
	}
}
