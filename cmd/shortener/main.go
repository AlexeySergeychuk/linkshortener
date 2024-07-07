package main

import (
	"net/http"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/handlers"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/repository"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/shortener"
)

func main() {
	mapStorage := make(map[string]string)
	repo := repository.NewRepo(mapStorage)
	shortener := shortener.NewShortener(repo)
	handler := handlers.NewHandler(shortener)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.CreateLinkHandler)
	mux.HandleFunc("/{id}", handler.GetLinkHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
