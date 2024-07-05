package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	APIKey      string `json:"metals_api_key"`
	REDISPW     string `json:"redis_pw"`
	PriceURL    string `json:"PriceURL"`
	NewsFeedURL string `json:"NewsFeedURL"`
}

func loadConfig(filename string) (Config, error) {
	var config Config

	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
