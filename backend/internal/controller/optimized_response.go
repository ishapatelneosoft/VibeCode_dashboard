package controller

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OptimizedJSONResponse provides faster JSON serialization for common responses
type OptimizedJSONResponse struct{}

// FastSuccess returns a success response with optimized JSON marshaling
func (ojr *OptimizedJSONResponse) FastSuccess(ctx *gin.Context, data interface{}) {
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)

	// Use buffer pool for JSON encoding
	buf := getBuffer()
	defer putBuffer(buf)

	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false) // Faster without HTML escaping

	if err := encoder.Encode(Response{
		Success: true,
		Data:    data,
	}); err != nil {
		// Fallback to standard JSON
		ctx.JSON(http.StatusOK, Response{
			Success: true,
			Data:    data,
		})
		return
	}

	ctx.Writer.Write(buf.Bytes())
}

// FastError returns an error response with optimized JSON marshaling
func (ojr *OptimizedJSONResponse) FastError(ctx *gin.Context, statusCode int, err error) {
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteHeader(statusCode)

	buf := getBuffer()
	defer putBuffer(buf)

	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(Response{
		Success: false,
		Error:   err.Error(),
	}); err != nil {
		// Fallback to standard JSON
		ctx.JSON(statusCode, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.Writer.Write(buf.Bytes())
}

// Precomputed responses for common endpoints
var (
	healthResponse       []byte
	notFoundResponse     []byte
	unauthorizedResponse []byte
)

func init() {
	// Precompute common responses
	healthResponse, _ = json.Marshal(Response{
		Success: true,
		Data:    gin.H{"status": "ok"},
	})

	notFoundResponse, _ = json.Marshal(Response{
		Success: false,
		Error:   "Not Found",
	})

	unauthorizedResponse, _ = json.Marshal(Response{
		Success: false,
		Error:   "Unauthorized",
	})
}

// FastHealthResponse returns precomputed health check response
func FastHealthResponse(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(healthResponse)
}

// FastNotFoundResponse returns precomputed 404 response
func FastNotFoundResponse(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusNotFound)
	ctx.Writer.Write(notFoundResponse)
}

// FastUnauthorizedResponse returns precomputed 401 response
func FastUnauthorizedResponse(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusUnauthorized)
	ctx.Writer.Write(unauthorizedResponse)
}

// Buffer pool for JSON encoding
var bufferPool = make(chan *bytes.Buffer, 100)

func getBuffer() *bytes.Buffer {
	select {
	case buf := <-bufferPool:
		buf.Reset()
		return buf
	default:
		return &bytes.Buffer{}
	}
}

func putBuffer(buf *bytes.Buffer) {
	select {
	case bufferPool <- buf:
		// Buffer returned to pool
	default:
		// Pool full, discard buffer
	}
}

// StructFieldOptimization provides hints for JSON field ordering
type StructFieldOptimization struct {
	// Fields ordered by frequency of access (most frequent first)
	FieldOrder []string
}

// OptimizeStruct reorders struct fields for better cache locality
// This is a compile-time optimization suggestion
func OptimizeStruct(s interface{}) {
	// In practice, this would use reflection to reorder fields
	// For now, it's a documentation of the optimization pattern
	_ = s
}
