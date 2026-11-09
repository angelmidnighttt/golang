package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger ghi log mỗi request bằng zerolog (thay cho logger mặc định của gin).
// Level theo status: 5xx -> error, 4xx -> warn, còn lại -> info.
// Cách ghi request/response body: xem describeBody trong body.go.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			path += "?" + maskSensitive(raw)
		}

		// Phải đọc request body trước c.Next(), trước khi handler đọc mất.
		requestHead := peekBody(c.Request)

		// Bọc writer để chép lại response body khi handler ghi ra.
		writer := &responseBodyWriter{ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		// Mô tả body sau c.Next() để lấy được kích thước file chính xác
		// nếu handler đã parse multipart form (c.FormFile, c.MultipartForm...).
		requestBody := describeBody(requestHead, c.Request.Header.Get("Content-Type"),
			c.Request.ContentLength, c.Request.MultipartForm)
		responseBody := describeBody(writer.body.Bytes(), writer.Header().Get("Content-Type"),
			int64(writer.Size()), nil)

		status := c.Writer.Status()

		var event *zerolog.Event
		switch {
		case status >= 500:
			event = log.Error()
		case status >= 400:
			event = log.Warn()
		default:
			event = log.Info()
		}

		if len(c.Errors) > 0 {
			event = event.Str("errors", c.Errors.String())
		}

		event.
			Str("method", c.Request.Method).
			Str("path", path).
			Str("request_body", requestBody).
			Int("status", status).
			Str("response_body", responseBody).
			Dur("latency", time.Since(start)).
			Str("ip", c.ClientIP()).
			Msg("request")
	}
}
