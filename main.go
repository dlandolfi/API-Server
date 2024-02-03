package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// todo; upgrade to structured logging
// todo: set up routes.go
// todo: set up tests
// todo: consider json -> toml
// todo: make port flag

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

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/api/v1/testpublic", testPublic)
	mux.Handle("/api/v1/testprivate", authMiddleware(http.HandlerFunc(testPrivate)))
	mux.Handle("/api/v1/testauthenticated", authMiddleware(http.HandlerFunc(testAuthenticated)))

	fmt.Println("Server is running on port:", port)
	if err := http.ListenAndServe(":"+port, securityHeadersMiddleware(mux)); err != nil {
		return err
	}
	return nil
}
