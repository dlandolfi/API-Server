package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// todo; upgrade to structured logging
// todo: set up routes.go
// todo: set up tests

var localMode = true
var logfile = "./server.log"

type Config struct {
	HRMS        HRMS        `json:"HRMS"`
	SSOProvider SSOProvider `json:"SSOProvider"`
}

type HRMS struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type SSOProvider struct {
	UserInfoURL string `json:"userInfoUrl"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	port := "3002"

	// Set up logging
	// If localMode, use standard out instead of a file
	var logFile *os.File
	var err error
	if localMode {
		logFile = os.Stdout
	} else {

		logFile, err = os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/apix/test", apixTest)
	mux.Handle("/api/test", authMiddleware(http.HandlerFunc(apiTest)))
	mux.Handle("/api/ps", authMiddleware(http.HandlerFunc(makeAuthenticatedCall)))

	fmt.Println("Server is running on port:", port)
	if err := http.ListenAndServe(":"+port, securityHeadersMiddleware(mux)); err != nil {
		return err
	}
	return nil
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
