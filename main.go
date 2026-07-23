package main

import "github.com/gin-gonic/gin"


func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/user/:user_id", func(c *gin.Context){
		userId := c.Param("user_id")
		c.JSON(200,gin.H{
			"message": userId,
		})
	})

	r.GET("/user",func (c *gin.Context){
		id := c.Query("id")
		name := c.Query("name")
		c.JSON(200,gin.H{
			"id":id,
			"name":name,
		})
	})

	r.Run(":8080")
}

