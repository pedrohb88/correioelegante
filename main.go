package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Message struct {
	ID   int    `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

func main() {
	// Open the SQLite database
	var err error
	db, err = sql.Open("sqlite3", "messages.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the messages table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			from_user TEXT,
			to_user TEXT,
			text TEXT,
			read BOOL DEFAULT false
		)
	`)
	if err != nil {
		log.Fatal("failed to create database: ", err)
	}

	// Define HTTP routes
	http.Handle("/", http.FileServer(http.Dir("."))) // Serve static files from the current directory
	http.HandleFunc("/messages", handle)
	http.HandleFunc("/admin", serveAdminPage)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveAdminPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("password") != "senhajunina" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	http.ServeFile(w, r, "admin.html")
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleGetMessages(w, r)
	} else if r.Method == http.MethodPost {
		handleMessage(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	var message Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO messages (from_user, to_user, text) VALUES (?, ?, ?)", message.From, message.To, message.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleGetMessages(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, from_user, to_user, text FROM messages WHERE read != true")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	messages := []Message{}

	var (
		updatePlaceholders []string
		updateArgs         []interface{}
	)
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.From, &msg.To, &msg.Text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)

		updatePlaceholders = append(updatePlaceholders, "?")
		updateArgs = append(updateArgs, msg.ID)
	}

	if len(updateArgs) > 0 {
		query := fmt.Sprintf("UPDATE messages SET read = true WHERE id IN(%s)", strings.Join(updatePlaceholders, ","))

		_, updateErr := db.Exec(query, updateArgs...)
		if updateErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
