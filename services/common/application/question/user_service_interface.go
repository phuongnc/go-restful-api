package question

import (
	"github.com/gin-gonic/gin"
)

type UserServiceHandler interface {
	UpdateUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
}
