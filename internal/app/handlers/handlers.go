package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/shortener"
)

// Создаем короткую ссылку
func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
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
	shortLink := shortener.MakeShortLink(string(body))

	log.Printf("ShortLink: %s", shortLink)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, shortLink)
}

// Вовращаем полную ссылкуу
func GetLinkHandler(w http.ResponseWriter, r *http.Request) {
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