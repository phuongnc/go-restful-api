package question

import (
	"github.com/gin-gonic/gin"
)

type AuthServiceHandler interface {
	Login(ctx *gin.Context)
}
