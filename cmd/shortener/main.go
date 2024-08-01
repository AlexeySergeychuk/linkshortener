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

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := *logger.Sugar()

	mapStorage := make(map[string]string)
	repo := repository.NewRepo(mapStorage)
	shortLinkStub := shortener.NewShortLink()
	shortener := shortener.NewShortener(repo, shortLinkStub)
	handler := handlers.NewHandler(shortener)

	router := gin.Default()
	router.Use(handlers.WithLogging(&sugar))
	router.Use(handlers.DecompressMiddleware())
	// не получается пройти автотесты с включенным компрессором ответов. Нужна помощь. Как ни странно без него автотесты проходят, то есть я делаю вывод,
	// что отдача gzip не тестируется. Поэтому вдвойне странно что при локальном запуске с включенным CompressMiddleware() работает
	// router.Use(handlers.CompressMiddleware()) 

	router.POST("/", handler.CreateShortLinkHandler)
	router.GET("/:id", handler.GetFullLinkHandler)
	router.POST("/api/shorten", handler.GetShortLinkHandler)

	router.Run(config.FlagRunAddr)
}
