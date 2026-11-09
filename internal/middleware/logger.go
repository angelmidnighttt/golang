package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger ghi log mỗi request bằng zerolog (thay cho logger mặc định của gin).
// Level theo status: 5xx -> error, 4xx -> warn, còn lại -> info.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			path += "?" + raw
		}

		c.Next()

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
			Int("status", status).
			Dur("latency", time.Since(start)).
			Str("ip", c.ClientIP()).
			Msg("request")
	}
}
