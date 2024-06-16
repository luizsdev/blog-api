package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	dbStr = "postgresql://postgres:postgres@localhost:5432/go?sslmode=disable"
)

type UserData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/user", createUserHandler)

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)

}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var user UserData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	db, err := sql.Open("postgres", dbStr)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	_, err = db.Query("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, hashedPassword)

	if err != nil {
		http.Error(w, "error inserting user", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)

	defer db.Close()
}
