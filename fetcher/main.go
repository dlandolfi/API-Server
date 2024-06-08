package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func fetchPrice() ([]byte, error) {
	resp, err := http.Get("https://api.metals.dev/v1/metal/spot?api_key=1XXURGAUJCZZAWTFJPHB808TFJPHB&metal=gold&currency=USD")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch gold price: status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func main() {
	// ExampleClient()
	response, err := fetchPrice()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(response))
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
