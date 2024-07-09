package main

import (
	"net/http"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/handlers"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/repository"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/shortener"
)

type HandlerService interface {
	CreateLinkHandler(w http.ResponseWriter, r *http.Request)
	GetLinkHandler(w http.ResponseWriter, r *http.Request)
}

func main() {
	mapStorage := make(map[string]string)
	repo := repository.NewRepo(mapStorage)
	shortLinkStub := shortener.NewShortLink()
	shortener := shortener.NewShortener(repo, shortLinkStub)
	handler := handlers.NewHandler(shortener)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.CreateLinkHandler)
	mux.HandleFunc("/{id}", handler.GetLinkHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
