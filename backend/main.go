package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/m/v2/sqlc/api"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// todo; upgrade to structured logging
// todo: set up routes.go
// todo: set up tests
// todo: consider json -> toml

var localMode = true
var logfile = "./server.log"

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

	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/api/v1/testpublic", testPublic)
	r.Handle("/api/v1/testprivate", authMiddleware(http.HandlerFunc(testPrivate)))
	r.Handle("/api/v1/testauthenticated", authMiddleware(http.HandlerFunc(testAuthenticated)))
	r.HandleFunc("/api/v1/getuser", getUserHandler) // /getuser?id=n
	r.HandleFunc("/api/v1/insertuser", createUserInDb)

	fmt.Println("Server is running on port:", port)
	if err := http.ListenAndServe(":"+port, securityHeadersMiddleware(r)); err != nil {
		return err
	}
	return nil
}
func dbConnect(ctx context.Context) (*api.Queries, error) {
	connStr := "postgres://user:passwords@postgres:5432/API_DB?sslmode=disable"
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	queries := api.New(conn)
	return queries, nil
}
