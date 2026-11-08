package v1handler

import (
	"net/http"

	"angelmidnighttt/golang/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandlerV1 struct{}

func NewUserHandler() *UserHandlerV1 {
	return &UserHandlerV1{}
}

func (u *UserHandlerV1) GetUsersV1(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{
		"message":"List All Users",
	})
}

func (u *UserHandlerV1) GetUserV1(c *gin.Context) {
	v := utils.NewValidator()

	id, ok := utils.ParsePositiveInt(c.Param("id"))
	v.Check(ok, "id", "id must be a positive integer")

	if !v.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": id,
	})
}
