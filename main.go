package main

import (
	"context"
	"fmt"
	"reflect"
)

func main() {
	redisClient, err := getRedisClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	//Basic data sets for strings and integers
	//Follow the convention we have domain:id:param
	redisClient.Set(ctx, "athletes:15:name", "Jhonny", 0)
	redisClient.Set(ctx, "athletes:15:weight", 82.45, 0)
	redisClient.Set(ctx, "athletes:15:age", 25, 0)
	name, _ := redisClient.Get(ctx, "athletes:15:name").Result()
	age, _ := redisClient.Get(ctx, "athletes:15:age").Int()
	weight, _ := redisClient.Get(ctx, "athletes:15:weight").Result()

	fmt.Printf("Athlet with id:15 name: %v (type: %v)\n", name, reflect.TypeOf(name))
	fmt.Printf("Athlet with id:15 age: %v (type: %v)\n", age, reflect.TypeOf(age))
	fmt.Printf("Athlet with id:15 weight: %v (type: %v)\n", weight, reflect.TypeOf(weight))

	//Inc functions
	redisClient.Incr(ctx, "athletes:15:age")
	age, _ = redisClient.Get(ctx, "athletes:15:age").Int()
	fmt.Printf("---\nNew athlet age is: %v (type: %v)\n", age, reflect.TypeOf(age))
	redisClient.Decr(ctx, "athletes:15:age")
	redisClient.IncrBy(ctx, "athletes:15:age", 10)
	age, _ = redisClient.Get(ctx, "athletes:15:age").Int()
	fmt.Printf("Updated athlet age: %v (type: %v)\n", age, reflect.TypeOf(age))

	//finding by keys
	keys, _ := redisClient.Keys(ctx, "athletes:*").Result()
	fmt.Println("---\nkeys for pattern: 'athletes:*'")
	for i, v := range keys {
		fmt.Printf("%v. key: %v\n", i+1, v)
	}
}
