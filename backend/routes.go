package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
)

func setupRoutes(r *mux.Router) {
	// Defining routes
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/api/v1/testpublic", testPublic).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/api/v1/testprivate", authMiddleware(http.HandlerFunc(testPrivate)))
	r.Handle("/api/v1/testauthenticated", authMiddleware(http.HandlerFunc(testAuthenticated)))
	r.HandleFunc("/api/v1/getuser", getUserHandler).Methods(http.MethodGet, http.MethodOptions) // /getuser?id=n
	r.HandleFunc("/api/v1/insertuser", createUserInDb).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/getallusers", getAllUsersHandler).Methods(http.MethodGet, http.MethodOptions)

	h := handler.New(&handler.Config{
		Schema: &schema,
	})
	r.Handle("/api/graphql", h)

	// Applying middlewares
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(securityHeadersMiddleware)
	r.Use(noCacheHeader)
}
