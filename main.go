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
		userHandlerv1 := v1handler.NewUserHandler()
		v1.GET("/v1/users", userHandlerv1.GetUsersV1)
		v1.GET("/v1/users/:id", userHandlerv1.GetUserV1)
	}

	

	r.Run(":8080")
}

