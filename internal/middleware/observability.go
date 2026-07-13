package middleware

import (
	"net/http"
	"time"

	"github.com/chai-rs/simple-bookstore/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
	UserIDKey       = "user_id"
)

// RequestID propagates an incoming request id or creates one for the request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set(RequestIDHeader, requestID)
		c.Next()
	}
}

// StructuredRequestLogger logs request completion with correlation fields.
func StructuredRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := requestIDFromGin(c)
		spanContext := trace.SpanContextFromContext(c.Request.Context())

		context := log.With()
		if requestID != "" {
			context = context.Str(RequestIDKey, requestID)
		}
		if spanContext.IsValid() {
			context = context.
				Str("trace_id", spanContext.TraceID().String()).
				Str("span_id", spanContext.SpanID().String())
		}

		requestLogger := context.Logger()
		c.Request = c.Request.WithContext(requestLogger.WithContext(c.Request.Context()))

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		event := requestLogEvent(status, requestLogger).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("route", route).
			Int("status", status).
			Dur("latency", latency).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent())

		if requestID != "" {
			event = event.Str(RequestIDKey, requestID)
		}
		if userID := userIDFromGin(c); userID != "" {
			event = event.Str(UserIDKey, userID)
		}
		if spanContext.IsValid() {
			event = event.
				Str("trace_id", spanContext.TraceID().String()).
				Str("span_id", spanContext.SpanID().String())
		}
		if len(c.Errors) > 0 {
			event = event.Str("gin_errors", c.Errors.String())
		}

		event.Msg("http request completed")
	}
}

// StructuredRecovery recovers panics and writes them through the structured logger.
func StructuredRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Ctx(c.Request.Context()).Error().
			Interface("panic", recovered).
			Msg("panic recovered")

		utils.ResponseErrorWithStatus(c, http.StatusInternalServerError, "internal server error")
	})
}

func requestIDFromGin(c *gin.Context) string {
	value, ok := c.Get(RequestIDKey)
	if !ok {
		return ""
	}

	requestID, ok := value.(string)
	if !ok {
		return ""
	}
	return requestID
}

func userIDFromGin(c *gin.Context) string {
	value, ok := c.Get(UserIDKey)
	if !ok {
		return ""
	}

	userID, ok := value.(string)
	if !ok {
		return ""
	}
	return userID
}

func requestLogEvent(status int, logger zerolog.Logger) *zerolog.Event {
	switch {
	case status >= 500:
		return logger.Error()
	case status >= 400:
		return logger.Warn()
	default:
		return logger.Info()
	}
}
