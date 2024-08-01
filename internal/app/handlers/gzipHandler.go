package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CompressMiddleware добавляет поддержку сжатия данных в ответах
func CompressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, поддерживает ли клиент сжатие данных
		acceptEncoding := c.Request.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		// Создаем новый compressWriter, если клиент поддерживает gzip
		var cw *compressWriter
		if supportsGzip {
			cw = newCompressWriter(c.Writer)
			c.Writer = cw
			defer cw.Close()
		}

		// Проверяем, отправлены ли сжатые данные от клиента
		contentEncoding := c.Request.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			cr, err := newCompressReader(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decompress request body"})
				c.Abort()
				return
			}
			c.Request.Body = cr
			defer cr.Close()
		}

		c.Next()

		if cw != nil {
			cw.Flush()
		}
	}
}

// DecompressMiddleware добавляет поддержку распаковки данных из запросов
func DecompressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, если запрос сжат
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			reader, err := newCompressReader(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gzip data"})
				c.Abort()
				return
			}

			c.Request.Body = reader
		}
		c.Next()
	}
}
