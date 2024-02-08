package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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

	makeDbConnection()

	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/api/v1/testpublic", testPublic)
	r.Handle("/api/v1/testprivate", authMiddleware(http.HandlerFunc(testPrivate)))
	r.Handle("/api/v1/testauthenticated", authMiddleware(http.HandlerFunc(testAuthenticated)))

	fmt.Println("Server is running on port:", port)
	if err := http.ListenAndServe(":"+port, securityHeadersMiddleware(r)); err != nil {
		return err
	}
	return nil
}

func makeDbConnection() {
	// PostgreSQL connection parameters
	connStr := "postgres://user:passwords@postgres:5432/API_DB?sslmode=disable"

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to PostgreSQL!")
}
