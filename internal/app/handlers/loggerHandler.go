package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// WithLogging добавляет логирование для всех запросов
func WithLogging(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// время начала обработки запроса
		start := time.Now()

		// передаем выполнение следующему обработчику
		c.Next()

		// время завершения обработки запроса
		duration := time.Since(start)

		// метод и URI запроса
		method := c.Request.Method
		uri := c.Request.RequestURI

		// статус и размер ответа
		status := c.Writer.Status()
		responseSize := c.Writer.Size()

		// логирование с использованием zap
		logger.Infow("request completed",
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", status,
			"responseSize", responseSize,
		)
	}
}