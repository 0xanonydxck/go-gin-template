package utils

import (
	"net/http"

	errs "github.com/chai-rs/simple-bookstore/internal/error"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Response represents the response structure.
type Response struct {
	Success bool `json:"success"`
	Error   any  `json:"error,omitempty"`
	Result  any  `json:"result,omitempty"`
}

// ResponseOk sends a success response.
func ResponseOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Result:  data,
	})
}

// ResponseCreated sends a created response.
func ResponseCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Result:  data,
	})
}

// ResponseErrorWithStatus sends an error response with a specific status.
func ResponseErrorWithStatus(c *gin.Context, status int, errorMessage string) {
	responseErrorWithStatus(c, status, errorMessage, nil)
}

// ResponseError sends an error response.
func ResponseError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errs.AppError:
		responseErrorWithStatus(c, e.Code, e.Message, err)
	default:
		responseErrorWithStatus(c, http.StatusInternalServerError, "internal server error", err)
	}
}

func responseErrorWithStatus(c *gin.Context, status int, errorMessage string, err error) {
	recordResponseError(c, status, errorMessage, err)
	c.JSON(status, Response{Success: false, Error: errorMessage})
}

func recordResponseError(c *gin.Context, status int, message string, err error) {
	span := trace.SpanFromContext(c.Request.Context())
	if !span.IsRecording() {
		return
	}

	span.SetAttributes(
		attribute.Int("http.response.status_code", status),
		attribute.String("error.message", message),
	)
	if err != nil {
		span.RecordError(err)
	}
	if status >= http.StatusInternalServerError {
		span.SetStatus(codes.Error, message)
	}
}
