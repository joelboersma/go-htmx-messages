package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const keyServerAddr = "serverAddr"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hasFirst := r.URL.Query().Has("first")
		first := r.URL.Query().Get("first")
		hasSecond := r.URL.Query().Has("second")
		second := r.URL.Query().Get("second")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("could not read body: %s\n", err)
		}

		res := fmt.Sprintf("first(%t)=%s, second(%t)=%s\nbody:\n%s\n",
			hasFirst, first,
			hasSecond, second,
			body,
		)

		io.WriteString(w, res)
	})

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		myName := r.PostFormValue("myName")
		if myName == "" {
			myName = "HTTP"
		}
		io.WriteString(w, fmt.Sprintf("Hello %s!\n", myName))
	})

	fmt.Println("Running server on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error listening for server: %s\n", err)
	}
}
