package main

import (
	"fmt"
	"log"
	"net/http"
)

// Security Headers
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r)
	})
}

func noCacheHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")

		next.ServeHTTP(w, r)
	})
}

// Authentication Middleware
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isAuthenticated(r) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func isAuthenticated(r *http.Request) bool {
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}
	url := config.SSOProvider.UserInfoURL
	authHeader := r.Header.Get("Authorization")

	if authHeader != "" {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("Authorization", authHeader)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode == 200 {
			return true
		}
	}
	return false
}
