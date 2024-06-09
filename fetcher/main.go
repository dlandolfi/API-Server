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
	// Load config
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Build new http client
	url := "https://api.metals.dev/v1/metal/spot"
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Add query params
	q := req.URL.Query()
	q.Add("api_key", config.APIKey)
	q.Add("metal", "gold")
	q.Add("currency", "USD")

	// Build new URL
	req.URL.RawQuery = q.Encode()

	// Execute request
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
	response, err := fetchPriceObject()
	if err != nil {
		fmt.Println(err)
	}
	err = rdb.Set(ctx, "priceObject", response, 0).Err()
	if err != nil {
		panic(err)
	}
}

func main() {
	var ctx = context.Background()
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: config.REDISPW,
		DB:       0, // use default DB
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
