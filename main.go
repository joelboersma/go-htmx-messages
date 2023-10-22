package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
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
		currentUser, err := getUser(db, 1)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't get sender: %s", err), http.StatusInternalServerError)
			return
		}
		otherUser, err := getUser(db, 2)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't get recipient: %s", err), http.StatusInternalServerError)
			return
		}

		messages, err := getMessages(db, currentUser.id, otherUser.id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to get messages: %s", err), http.StatusInternalServerError)
			return
		}

		displayMessages := []struct {
			SenderName string
			Content    string
		}{}
		for _, message := range messages {
			name := "Unknown"
			if message.sender_id == currentUser.id {
				name = currentUser.name
			} else if message.sender_id == otherUser.id {
				name = otherUser.name
			}

			displayMessage := struct {
				SenderName string
				Content    string
			}{
				name,
				message.content,
			}

			displayMessages = append(displayMessages, displayMessage)
		}

		tmpl, err := template.ParseFiles("templates/messages.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to parse template: %s", err), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, displayMessages)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to execute template: %s", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf8")
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
