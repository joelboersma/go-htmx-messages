package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const keyServerAddr = "serverAddr"

func main() {
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		htmlContent, err := os.ReadFile("templates/index.html")
		if err != nil {
			http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf8")
		w.Write(htmlContent)
	})

	mux.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		// TEMP: for testing
		userId1 := 1
		userId2 := 2

		messages, err := getMessages(db, userId1, userId2)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to read fetch messages: %s", err), http.StatusInternalServerError)
			return
		}

		fmt.Println(messages)

		htmlContent, err := os.ReadFile("templates/messages.html")
		if err != nil {
			http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf8")
		w.Write(htmlContent)
	})

	staticFs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))

	fmt.Println("Running server on port 8080")
	err = http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error listening for server: %s\n", err)
	}

	db.Close()
}
