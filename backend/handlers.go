package main

import (
	"api-server/sqlc/api"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

// Handlers
func getNewsFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rubyServerURL := "http://ruby_server:3000/api/newsfeed"

	resp, err := http.Get(rubyServerURL)
	if err != nil {
		http.Error(w, "Failed to fetch data from ruby_server", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from ruby_server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
func getPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: config.REDISPW,
		DB:       0, // use default DB
	})

	val, err := rdb.Get(ctx, "priceObject").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key does not exist")
		} else {
			panic(err)
		}
	}

	fmt.Fprint(w, val)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()
	queries, err := dbConnect(ctx)
	if err != nil {
		log.Println(err)
	}

	user, err := queries.GetUser(ctx, int32(id))
	if err != nil {
		http.Error(w, "An error has occured", http.StatusBadRequest)
		log.Println("queries.getuser", err)
		return
	}
	fmt.Println(user)

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		log.Println("Error: ", err)
	}
	result := string(jsonBytes)
	io.WriteString(w, result)
}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	queries, err := dbConnect(ctx)
	if err != nil {
		log.Println(err)
	}

	allUsers, err := queries.GetAllUsers(ctx)
	if err != nil {
		log.Println(err)
	}

	jsonData, err := json.Marshal(allUsers)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonData)
}

func createUserInDb(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	queries, err := dbConnect(ctx)
	if err != nil {
		log.Println(err)
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")

	insertedUser, err := queries.CreateUser(ctx, api.CreateUserParams{
		FirstName: pgtype.Text{String: firstName, Valid: true},
		LastName:  pgtype.Text{String: lastName, Valid: true},
		Email:     pgtype.Text{String: email, Valid: true},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Create user error:", err)
		return
	}

	log.Println("Inserted user:", insertedUser)
	io.WriteString(w, "Successfully created user")

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to API server.")
}

func testPublic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	log.Println("Public Test route hit!")
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
