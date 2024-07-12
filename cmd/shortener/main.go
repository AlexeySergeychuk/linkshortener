package main

import (
	"github.com/gin-gonic/gin"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/config"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/handlers"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/repository"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/shortener"
)

func main() {
	config.ParseFlags()

	mapStorage := make(map[string]string)
	repo := repository.NewRepo(mapStorage)
	shortLinkStub := shortener.NewShortLink()
	shortener := shortener.NewShortener(repo, shortLinkStub)
	handler := handlers.NewHandler(shortener)

	router := gin.Default()
	router.POST("/", handler.CreateLinkHandler)
	router.GET("/:id", handler.GetLinkHandler)

	router.Run(config.FlagRunAddr)
}
