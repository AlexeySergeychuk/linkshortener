package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/config"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/handlers"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/repository"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/shortener"
)

func main() {
	config.ParseFlags()

	// создаём предустановленный регистратор zap
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar := *logger.Sugar()

	mapStorage := make(map[string]string)
	repo := repository.NewRepo(mapStorage)
	shortLinkStub := shortener.NewShortLink()
	shortener := shortener.NewShortener(repo, shortLinkStub)
	handler := handlers.NewHandler(shortener)

	router := gin.Default()
	router.Use(handlers.WithLogging(&sugar))
	router.POST("/", handler.CreateLinkHandler)
	router.GET("/:id", handler.GetLinkHandler)

	router.Run(config.FlagRunAddr)
}
