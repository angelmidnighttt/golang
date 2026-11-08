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

type GetUserParamV1 struct {
	ID int `uri:"id" binding:"gt=0"`
}

func (u *UserHandlerV1) GetUserV1(c *gin.Context) {
	var params GetUserParamV1
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": params.ID,
	})
}
