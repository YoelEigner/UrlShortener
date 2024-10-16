package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the URL Shortener API!")
	})

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
