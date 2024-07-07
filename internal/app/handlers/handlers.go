package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/shortener"
)

type HandlerService interface {
	CreateLinkHandler(w http.ResponseWriter, r *http.Request)
	GetLinkHandler(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	shortener shortener.ShortenerService
}

func NewHandler(s shortener.ShortenerService) *Handler{
	return &Handler{shortener: s}
}

// Создаем короткую ссылку
func (h *Handler) CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	defer r.Body.Close()

	log.Printf("Body: %s", string(body))
	// Отправка тела запроса обратно в ответ
	shortLink := h.shortener.MakeShortLink(string(body))

	log.Printf("ShortLink: %s", shortLink)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, shortLink)
}

// Вовращаем полную ссылку
func (h *Handler) GetLinkHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Printf("Path: %s", path)

	fullLink := h.shortener.GetFullLink(path)
	if fullLink == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Printf("Full link: %s", fullLink)

	w.Header().Set("Location", fullLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}