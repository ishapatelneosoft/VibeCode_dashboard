package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SuccessResponse creates a successful response
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessResponseWithMessage creates a successful response with a message
func SuccessResponseWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse creates an error response
func ErrorResponse(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   err.Error(),
	})
}

// BadRequestResponse creates a 400 Bad Request response
func BadRequestResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusBadRequest, err)
}

// UnauthorizedResponse creates a 401 Unauthorized response
func UnauthorizedResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusUnauthorized, err)
}

// NotFoundResponse creates a 404 Not Found response
func NotFoundResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusNotFound, err)
}

// InternalServerErrorResponse creates a 500 Internal Server Error response
func InternalServerErrorResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusInternalServerError, err)
}

// ValidationErrorResponse creates a 422 Unprocessable Entity response for validation errors
func ValidationErrorResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusUnprocessableEntity, err)
}
