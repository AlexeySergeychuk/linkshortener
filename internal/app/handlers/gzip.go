package handlers

import (
	"bufio"
	"compress/gzip"
	"io"
	"net"
	"net/http"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w             http.ResponseWriter
	zw            *gzip.Writer
	size          int
	status        int
	headerWritten bool
	written       bool
}

// newCompressWriter создает новый compressWriter
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	if !c.headerWritten {
        c.w.Header().Set("Content-Encoding", "gzip")
        c.headerWritten = true
    }
    return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if !c.headerWritten {
        if statusCode < 300 {
            c.w.Header().Set("Content-Encoding", "gzip")
        }
        c.w.WriteHeader(statusCode)
        c.headerWritten = true
    }
}

// WriteHeaderNow записывает заголовки немедленно
func (c *compressWriter) WriteHeaderNow() {
	if !c.headerWritten {
		c.w.Header().Set("Content-Encoding", "gzip")
		c.w.WriteHeader(c.status)
		c.headerWritten = true
	}
}

// Flush отправляет все буферизованные данные в клиент и может использоваться
// для принудительной отправки данных, а также для поддержки постоянных соединений.
func (c *compressWriter) Flush() {
	if f, ok := c.w.(http.Flusher); ok {
		f.Flush()
	}
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
    if err := c.zw.Close(); err != nil {
        return err
    }
    return nil
}

// CloseNotify не поддерживается в gzip.Writer, поэтому возвращаем канал,
// который никогда не закрывается.
func (c *compressWriter) CloseNotify() <-chan bool {
	ch := make(chan bool)
	return ch
}

func (c *compressWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := c.w.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

func (c *compressWriter) Pusher() http.Pusher {
	if pusher, ok := c.w.(http.Pusher); ok {
		return pusher
	}
	return nil
}

func (c *compressWriter) Size() int {
	return c.size
}

func (c *compressWriter) Status() int {
	return c.status
}

func (c *compressWriter) WriteString(s string) (int, error) {
	return c.zw.Write([]byte(s))
}

func (c *compressWriter) Written() bool {
	return c.written
}

// newCompressReader создает новый compressReader
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// compressReader реализует интерфейс io.ReadCloser и позволяет декомпрессировать получаемые данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
