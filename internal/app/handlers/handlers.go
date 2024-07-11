package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShortenerService interface {
	MakeShortLink(link string) string
	GetFullLink(shortLink string) string
}

type Handler struct {
	shortener ShortenerService
}

func NewHandler(s ShortenerService) *Handler {
	return &Handler{shortener: s}
}

// Создаем короткую ссылку
func (h *Handler) CreateLinkHandler(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	defer c.Request.Body.Close()

	log.Printf("Body: %s", string(body))
	// Отправка тела запроса обратно в ответ
	shortLink := h.shortener.MakeShortLink(string(body))
	log.Printf("ShortLink: %s", shortLink)

	c.String(http.StatusCreated, shortLink)
}

// Вовращаем полную ссылку
func (h *Handler) GetLinkHandler(c *gin.Context) {
	shortLink := c.Param("id")

	log.Printf("shortLink in request: %s", shortLink)
	fullLink := h.shortener.GetFullLink("/" + shortLink)
	if fullLink == "" {
		c.String(http.StatusNotFound, "Link not found")
		return
	}
	log.Printf("Fulllink: %s", fullLink)

	c.Redirect(http.StatusTemporaryRedirect, fullLink)
}