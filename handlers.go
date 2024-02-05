package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// Handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("log homehandler")
	io.WriteString(w, "Hello World")
}

func testPublic(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Public route successful")
}

func testPrivate(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Protected route successful")
}

func testAuthenticated(w http.ResponseWriter, r *http.Request) {
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}
	url := config.HRMS.URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	req.Header.Add("Authorization", config.HRMS.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	io.WriteString(w, string(body))
}
