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

func fetchPriceObject(apiKey string) (string, error) {
	url := "https://api.metals.dev/v1/metal/spot"
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	q := req.URL.Query()
	q.Add("api_key", apiKey)
	q.Add("metal", "gold")
	q.Add("currency", "USD")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}
	return string(body), nil
}

func fetchAndStore(rdb *redis.Client, ctx context.Context, apiKey string) error {
	response, err := fetchPriceObject(apiKey)
	if err != nil {
		return fmt.Errorf("fetching price object: %w", err)
	}

	err = rdb.Set(ctx, "priceObject", response, 0).Err()
	if err != nil {
		return fmt.Errorf("storing in Redis: %w", err)
	}

	return nil
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: config.REDISPW,
		DB:       0, // use default DB
	})
	defer rdb.Close()

	ctx := context.Background()

	// Run the fetch and store function immediately
	go func() {
		if err := fetchAndStore(rdb, ctx, config.APIKey); err != nil {
			log.Printf("Error in initial fetch and store: %v", err)
		}
	}()

	// Calculate the duration until the next 8am
	now := time.Now()
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
	if now.After(nextTick) {
		nextTick = nextTick.Add(24 * time.Hour)
	}
	durationUntilNextTick := nextTick.Sub(now)

	time.AfterFunc(durationUntilNextTick, func() {
		if err := fetchAndStore(rdb, ctx, config.APIKey); err != nil {
			log.Printf("Error in fetch and store at 8am: %v", err)
		}

		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			if err := fetchAndStore(rdb, ctx, config.APIKey); err != nil {
				log.Printf("Error in fetch and store on ticker: %v", err)
			}
		}
	})

	select {}
}
