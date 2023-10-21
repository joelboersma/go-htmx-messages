package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World!\n")
	})

	fmt.Println("Running server on port 8080")
	http.ListenAndServe(":8080", nil)
}
