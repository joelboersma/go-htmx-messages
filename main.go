package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

const keyServerAddr = "serverAddr"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		htmlContent, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf8")
		w.Write(htmlContent)
	})

	fmt.Println("Running server on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error listening for server: %s\n", err)
	}
}
