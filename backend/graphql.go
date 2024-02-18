package main

import (
	"database/sql"
	"log"

	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
)

// User represents the User type in GraphQL
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.Int},
		"first_name": &graphql.Field{Type: graphql.String},
		"last_name":  &graphql.Field{Type: graphql.String},
		"email":      &graphql.Field{Type: graphql.String},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"users": &graphql.Field{
			Type: graphql.NewList(userType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db, err := sql.Open("postgres", "postgres://user:passwords@postgres:5432/API_DB?sslmode=disable")
				if err != nil {
					log.Fatal(err)
				}
				defer db.Close()

				rows, err := db.Query("SELECT id, first_name, last_name, email FROM users")
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()

				var users []User
				for rows.Next() {
					var user User
					err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
					if err != nil {
						log.Fatal(err)
					}
					users = append(users, user)
				}

				return users, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})
