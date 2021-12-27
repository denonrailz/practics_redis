package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-redis/redis/v8"
	. "github.com/logrusorgru/aurora"
	"math/rand"
	"strconv"
	"time"
)

const (
	maxResponseMsec = 1000
	usersCount      = 50
	requestsCount   = 200
	cacheTTL        = 200 //seconds
)

var rnd *rand.Rand

type User struct {
	ID    int    `json:"id" fake:"{number:1,1000}"`
	Name  string `json:"name" fake:"{firstname}"`
	Email string `json:"email" fake:"{email}"`
	Phone string `json:"phone" fake:"{phone}"`
}

type cache struct {
	ttl    time.Duration
	client *redis.Client
	ctx    context.Context
}

func init() {
	src := rand.NewSource(time.Now().UnixNano())
	rnd = rand.New(src)
}

//userRequest return same data on same userID
func userRequest(userID string) ([]byte, error) {
	var user User
	seedID, _ := strconv.ParseInt(userID, 10, 64)
	gofakeit.Seed(seedID)
	//random sleep
	time.Sleep(time.Duration(rnd.Intn(maxResponseMsec)) * time.Millisecond)
	//filling up structure
	_ = gofakeit.Struct(&user)
	return json.Marshal(user)
}

func (ch *cache) GetUserFromCache(userId string) (User, error) {
	byteUser, err := ch.client.Get(ch.ctx, userId).Bytes()
	if err != nil {
		return User{}, err
	}
	cached := toJson(byteUser)
	return cached, nil
}

func main() {
	redisClient, err := getRedisClient()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to redis: success")

	cache := cache{
		ttl:    cacheTTL * time.Second,
		client: redisClient,
		ctx:    context.Background(),
	}

	for i := 0; i < requestsCount; i++ {
		//getting random user id
		usrID := strconv.Itoa(rnd.Intn(usersCount))

		//resulted struct for user
		var resultUser User

		//flag for checking using cache
		userFromCache := true

		startTime := time.Now()
		resultUser, err := cache.GetUserFromCache(usrID)
		if err != nil {
			userFromCache = false
			//making lagging request
			body, err := userRequest(usrID)
			if err != nil {
				fmt.Printf("error occurs %v/n", err)
			}
			//setting up cache
			cacheErr := cache.client.Set(cache.ctx, usrID, body, cache.ttl).Err()
			if cacheErr != nil {
				fmt.Printf("cache error occurs %v/n", cacheErr)
			}
			resultUser = toJson(body)
		}

		//colored logs
		lag := time.Now().Sub(startTime)
		coloredFlag := Red("REQUESTED")
		coloredTime := Red(lag)
		if userFromCache {
			coloredFlag = Green("CACHED")
			coloredTime = Green(coloredTime)
		}
		fmt.Printf("id: %v, lag: %v cache: %v response: %+v \n", Yellow(usrID), coloredTime, coloredFlag, resultUser)
	}

	//flush all cache
	cache.client.FlushAll(cache.ctx)
}

func toJson(val []byte) User {
	user := User{}
	err := json.Unmarshal(val, &user)
	if err != nil {
		panic(err)
	}
	return user
}
