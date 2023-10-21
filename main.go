package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
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
		io.WriteString(w, "Hello HTTP!\n")
	})

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	fmt.Println("Running server on port 8080")
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error listening for server: %s\n", err)
	}
}
