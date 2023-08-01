package http

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type ErrorRes struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type Gin struct {
	C *gin.Context
}

func init() {
	// Register custom validate methods
	govalidator.TagMap["required"] = govalidator.Validator(func(str string) bool {
		return len(str) > 0
	})
}

func (g *Gin) Response(httpCode int, success bool, message string, data interface{}, err error) {
	g.C.JSON(httpCode, gin.H{
		"success": success,
		"message": message,
		"data":    data,
		"error":   err,
	})
	return
}

func (g *Gin) BindToDto(obj interface{}) bool {
	if err := g.C.ShouldBind(obj); err != nil {
		g.Response(http.StatusBadRequest, false, err.Error(), nil, nil)
		return false
	}
	return true
}
