package main

import (
	"io"
	"log"
	"net/http"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/shortener"
)

// Создаем короткую ссылку
func createLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	log.Printf("Body: %s", string(body))
	// Отправка тела запроса обратно в ответ
	shortLink := shortener.MakeShortLink(string(body));

	log.Printf("ShortLink: %s", shortLink)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, shortLink)
}

// Вовращаем полную ссылку
func getLinkHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Printf("Path: %s", path)

	fullLink := shortener.GetFullLink(path)
	if fullLink == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Printf("Full link: %s", fullLink)

	w.Header().Set("Location", fullLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", createLinkHandler)
	mux.HandleFunc("/{id}", getLinkHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
