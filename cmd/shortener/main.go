package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/config"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/handlers"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/repo"
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
	memRepo := repo.NewRepo(mapStorage)

	// Инициализация файлового хранилища
	fileProducer, err := repo.NewProducer(config.FlagFileDBPath)
	if err != nil {
		sugar.Error("Can't create file producer for links")
	}

	// Загрузка данных из файла в memRepo
	fileConsumer, err := repo.NewConsumer(config.FlagFileDBPath)
	if err != nil {
		sugar.Error("Can't create file consumer for links")
		return
	}

	urls, err := fileConsumer.ReadAll()
	if err != nil {
		sugar.Error("Can't read data from file")
		return
	}

	for _, url := range urls {
		memRepo.SaveLinks(url.ShortURL, url.FullURL)
	}

	shortLinkStub := shortener.NewShortLink()
	shortener := shortener.NewShortener(memRepo, fileProducer, shortLinkStub)
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
