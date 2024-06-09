package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	HRMS        HRMS        `json:"HRMS"`
	SSOProvider SSOProvider `json:"SSOProvider"`
	REDISPW     string      `json:"redis_pw"`
}

type HRMS struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type SSOProvider struct {
	UserInfoURL string `json:"userInfoUrl"`
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
