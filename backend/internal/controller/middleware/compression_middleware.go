package middleware

import (
	"compress/gzip"
	"io"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// gzipWriterPool reduces allocation overhead for gzip writers
var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(io.Discard)
	},
}

// CompressionMiddleware provides GZIP compression for responses
func CompressionMiddleware(minSize int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip compression for certain content types
		if !shouldCompress(ctx) {
			ctx.Next()
			return
		}

		// Check if client accepts gzip
		acceptEncoding := ctx.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			ctx.Next()
			return
		}

		// Create gzip writer from pool
		gzWriter := gzipWriterPool.Get().(*gzip.Writer)
		defer gzipWriterPool.Put(gzWriter)

		// Reset writer for new response
		gzWriter.Reset(ctx.Writer)

		// Create custom response writer
		writer := &compressedResponseWriter{
			ResponseWriter: ctx.Writer,
			gzipWriter:     gzWriter,
			minSize:        minSize,
			buffer:         make([]byte, 0),
		}
		ctx.Writer = writer

		// Set headers
		ctx.Header("Content-Encoding", "gzip")
		ctx.Header("Vary", "Accept-Encoding")

		// Cleanup
		defer func() {
			writer.Flush()
			gzWriter.Close()
		}()

		ctx.Next()
	}
}

// shouldCompress determines if response should be compressed
func shouldCompress(ctx *gin.Context) bool {
	// Skip compression for already compressed content
	contentType := ctx.GetHeader("Content-Type")

	// Skip for certain content types
	skipTypes := []string{
		"image/",
		"video/",
		"audio/",
		"application/octet-stream",
		"application/pdf",
		"application/zip",
		"application/gzip",
	}

	for _, skipType := range skipTypes {
		if strings.HasPrefix(contentType, skipType) {
			return false
		}
	}

	return true
}

// compressedResponseWriter buffers response to determine if compression is beneficial
type compressedResponseWriter struct {
	gin.ResponseWriter
	gzipWriter *gzip.Writer
	minSize    int
	buffer     []byte
	statusCode int
	compressed bool
}

func (w *compressedResponseWriter) Write(data []byte) (int, error) {
	// If already decided to compress, write compressed
	if w.compressed {
		return w.gzipWriter.Write(data)
	}

	// Buffer data to decide
	w.buffer = append(w.buffer, data...)

	// If buffer exceeds min size, compress everything
	if len(w.buffer) >= w.minSize {
		w.compressed = true
		w.ResponseWriter.Header().Set("Content-Encoding", "gzip")

		// Write buffered data compressed
		n, err := w.gzipWriter.Write(w.buffer)
		if err != nil {
			return n, err
		}

		// Clear buffer
		w.buffer = nil
		return len(data), nil
	}

	return len(data), nil
}

func (w *compressedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *compressedResponseWriter) Flush() {
	if w.compressed {
		w.gzipWriter.Flush()
	} else if len(w.buffer) > 0 {
		// Write uncompressed buffer
		w.ResponseWriter.Write(w.buffer)
	}
	w.ResponseWriter.Flush()
}

func (w *compressedResponseWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

// FastCompressionMiddleware is a simpler version that always compresses if beneficial
func FastCompressionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if client accepts gzip
		acceptEncoding := ctx.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			ctx.Next()
			return
		}

		// Skip for small responses or certain content types
		contentType := ctx.GetHeader("Content-Type")
		if strings.HasPrefix(contentType, "image/") ||
			strings.HasPrefix(contentType, "video/") ||
			strings.HasPrefix(contentType, "audio/") {
			ctx.Next()
			return
		}

		// Create gzip writer
		gzWriter := gzip.NewWriter(ctx.Writer)
		defer gzWriter.Close()

		// Set headers
		ctx.Header("Content-Encoding", "gzip")
		ctx.Header("Vary", "Accept-Encoding")

		// Create custom writer
		writer := &simpleCompressedWriter{
			ResponseWriter: ctx.Writer,
			gzipWriter:     gzWriter,
		}
		ctx.Writer = writer

		ctx.Next()
	}
}

// simpleCompressedWriter is a simple wrapper for gzip compression
type simpleCompressedWriter struct {
	gin.ResponseWriter
	gzipWriter *gzip.Writer
}

func (w *simpleCompressedWriter) Write(data []byte) (int, error) {
	return w.gzipWriter.Write(data)
}

func (w *simpleCompressedWriter) WriteString(s string) (int, error) {
	return w.gzipWriter.Write([]byte(s))
}
