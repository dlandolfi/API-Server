package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func fetchPriceObject() (string, error) {
	url := "https://api.metals.dev/v1/metal/spot"
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("api_key", "1XXURGAUJCZZAWTFJPHB808TFJPHB")
	q.Add("metal", "gold")
	q.Add("currency", "USD")

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	priceObject := string(body)
	return priceObject, nil
}

func fetchAndStore(rdb *redis.Client, ctx context.Context) {
	// response, err := fetchPriceObject()
	response, err := fakeFetch()
	if err != nil {
		fmt.Println(err)
	}
	err = rdb.Set(ctx, "priceObject", response, 0).Err()
	if err != nil {
		panic(err)
	}

}

func fakeFetch() (string, error) {
	return "{stuffsss}", nil
}

func main() {
	// ExampleClient()
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	go fetchAndStore(rdb, ctx)
	// Calculate the duration until the next 8am
	now := time.Now()
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
	if now.After(nextTick) {
		// If it's already past 8am today, set the next tick for tomorrow
		nextTick = nextTick.Add(24 * time.Hour)
	}
	durationUntilNextTick := nextTick.Sub(now)

	// Wait until the next 8am
	time.AfterFunc(durationUntilNextTick, func() {
		// Run fetchPrice at the next 8am
		fetchAndStore(rdb, ctx)

		// Set up a ticker to run fetchPrice every 24 hours from the next 8am
		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			fetchAndStore(rdb, ctx)
		}
	})

	// Keep the main goroutine running
	select {}

}

func ExampleClient() {

	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
