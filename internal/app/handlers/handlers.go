package handlers

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
)

type Shortener interface {
	MakeShortLink(link string) string
	GetFullLink(shortLink string) string
}

type Handler struct {
	shortener Shortener
}

func NewHandler(s Shortener) *Handler {
	return &Handler{shortener: s}
}

// Создаем короткую ссылку
func (h *Handler) CreateShortLinkHandler(c *gin.Context) {

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
func (h *Handler) GetFullLinkHandler(c *gin.Context) {
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

func (h *Handler) GetShortLinkHandler(c *gin.Context) {

	var link URLRequest
	if err := easyjson.UnmarshalFromReader(c.Request.Body, &link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	}

	defer c.Request.Body.Close()

	if link.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передана пустая ссылка"})
		return
	}

	shortLink := h.shortener.MakeShortLink(link.URL)

	urlResponse := URLResponse{
		Result: shortLink,
	}

	response, err := easyjson.Marshal(urlResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(response)))
	c.Data(http.StatusCreated, "application/json", response)
}
