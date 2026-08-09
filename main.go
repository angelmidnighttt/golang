package main

import (
	v1handler "angelmidnighttt/golang/internal/api/v1/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
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

	r.Run(":8080")
}
