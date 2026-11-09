package main

import (
	v1handler "angelmidnighttt/golang/internal/api/v1/handler"
	"angelmidnighttt/golang/internal/logger"
	"angelmidnighttt/golang/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	logger.Init()

	// gin.New() thay cho gin.Default() để dùng logger zerolog thay vì logger mặc định của gin.
	r := gin.New()

	// Không tin header X-Forwarded-For do client gửi lên,
	// nếu không client có thể giả IP để né rate limit.
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal().Err(err).Msg("set trusted proxies")
	}

	r.Use(
		middleware.Logger(),         // đặt đầu tiên để log cả request bị 429 hoặc panic
		gin.Recovery(),              // bắt panic, trả 500
		middleware.RateLimit(5, 10), // mỗi IP: 5 req/giây, dồn tối đa 10
	)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			userHandler := v1handler.NewUserHandler()
			user.GET("/", userHandler.GetUsersV1)
			user.GET("/:id", userHandler.GetUserV1)
		}
	}

	addr := ":8080"
	log.Info().Str("addr", addr).Msg("starting server")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}
