package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
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

func fetchNewsFeed() (string, error) {
	url := "http://ruby_server:3000/api/newsfeed"
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
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

	newsResponse, err := fetchNewsFeed()
	if err != nil {
		return fmt.Errorf("fetching news response: %w", err)
	}

	err = rdb.Set(ctx, "newsResponse", newsResponse, 0).Err()
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

	const maxRetries = 5
	const baseDelay = time.Second * 5

	for i := 0; i < maxRetries; i++ {
		if err := fetchAndStore(rdb, ctx, config.APIKey); err != nil {
			log.Printf("Error in initial fetch and store (attempt %d/%d): %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(baseDelay)
			}
		} else {
			break
		}
	}

	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		if err := fetchAndStore(rdb, ctx, config.APIKey); err != nil {
			log.Printf("Error in cron fetch and store: %v", err)
		}
	})
	c.Start()
	defer c.Stop()

	select {}
}
