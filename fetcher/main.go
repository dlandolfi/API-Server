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

const (
	maxRetries  = 5
	baseDelay   = 5 * time.Second
	priceURL    = "https://api.metals.dev/v1/metal/spot"
	newsFeedURL = "http://ruby_server:3000/api/newsfeed"
	timeout     = 10 * time.Second
)

func fetchURL(url string, params map[string]string) (string, error) {
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
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

func fetchPriceObject(apiKey string) (string, error) {
	params := map[string]string{
		"api_key":  apiKey,
		"metal":    "gold",
		"currency": "USD",
	}
	return fetchURL(priceURL, params)
}

func fetchNewsFeed() (string, error) {
	return fetchURL(newsFeedURL, nil)
}

func storeData(rdb *redis.Client, ctx context.Context, key string, value string) error {
	return rdb.Set(ctx, key, value, 0).Err()
}

func fetchAndStore(rdb *redis.Client, ctx context.Context, apiKey string) error {
	chanSize := 2
	errors := make(chan error, chanSize) // Channel to collect errors
	defer close(errors)

	go func() {
		priceData, err := fetchPriceObject(apiKey)
		if err != nil {
			log.Printf("Error fetching price object: %v", err)
			errors <- err
			return
		}
		if err := storeData(rdb, ctx, "priceObject", priceData); err != nil {
			errors <- fmt.Errorf("storing price object in Redis: %w", err)
		} else {
			errors <- nil
		}
	}()

	go func() {
		newsData, err := fetchNewsFeed()
		if err != nil {
			errors <- fmt.Errorf("fetching news feed: %w", err)
			return
		}
		if err := storeData(rdb, ctx, "newsResponse", newsData); err != nil {
			errors <- fmt.Errorf("storing news response in Redis: %w", err)
		} else {
			errors <- nil
		}
	}()

	// Wait for both goroutines to complete
	var finalErr error
	for i := 0; i < chanSize; i++ {
		if err := <-errors; err != nil {
			finalErr = err
		}
	}

	return finalErr
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: config.REDISPW,
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

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
