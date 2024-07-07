package main

import (
	"net/http"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/handlers"
)


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.CreateLinkHandler)
	mux.HandleFunc("/{id}", handlers.GetLinkHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
