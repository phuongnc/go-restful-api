package http

import (
	"net/http"
	"smartkid/services/common/infra/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	BadRequest(c *gin.Context, err error, message ...string)
	Unauthorized(c *gin.Context, err error, message ...string)
	InternalServerError(c *gin.Context, err error)
	ResponseError(c *gin.Context, message string, err error, httpCode ...int)
	ResponseData(c *gin.Context, data interface{}, message string, httpCode ...int)
	BindToDto(c *gin.Context, obj interface{}) error
	GetPagingParams(c *gin.Context) (offset, limit int32, err error)
}

type baseHandler struct {
	logger.Logger
}

func NewHandler(log logger.Logger) Handler {
	return &baseHandler{
		Logger: log,
	}
}

func (handler *baseHandler) BadRequest(c *gin.Context, err error, message ...string) {
	returnMessage := ""
	if len(message) > 0 {
		returnMessage = message[0]
	}
	handler.ResponseError(c, returnMessage, err, http.StatusBadRequest)
}

func (handler *baseHandler) InternalServerError(c *gin.Context, err error) {
	handler.ResponseError(c, "", err, http.StatusInternalServerError)
}

func (handler *baseHandler) Unauthorized(c *gin.Context, err error, message ...string) {
	returnMessage := ""
	if len(message) > 0 {
		returnMessage = message[0]
	}
	handler.ResponseError(c, returnMessage, err, http.StatusUnauthorized)
}

func (handler *baseHandler) ResponseError(c *gin.Context, message string, err error, httpCode ...int) {
	returnHttpCode := http.StatusBadRequest
	if len(httpCode) > 0 {
		returnHttpCode = httpCode[0]
	}
	handler.response(c, returnHttpCode, false, message, nil, err)
}

func (handler *baseHandler) ResponseData(c *gin.Context, data interface{}, message string, httpCode ...int) {
	returnHttpCode := http.StatusOK
	if len(httpCode) > 0 {
		returnHttpCode = httpCode[0]
	}
	handler.response(c, returnHttpCode, true, message, data, nil)
}

func (handler *baseHandler) response(c *gin.Context, httpCode int, success bool, message string, data interface{}, err error) {
	c.JSON(httpCode, gin.H{
		"success": success,
		"message": message,
		"data":    data,
		"error":   err,
	})
}

func (handler *baseHandler) BindToDto(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBind(obj); err != nil {
		handler.Logger.Error(err)
		handler.BadRequest(c, err)
	}
	return nil
}

func (handler *baseHandler) GetPagingParams(c *gin.Context) (offset, limit int32, err error) {
	if c.Query("offset") != "" {
		offsetParse, err := strconv.ParseInt(c.Query("offset"), 10, 32)
		if err != nil {
			return 0, 0, err
		}
		offset = int32(offsetParse)
	} else {
		offset = 0
	}
	if c.Query("limit") != "" {
		limitParse, err := strconv.ParseInt(c.Query("limit"), 10, 32)
		if err != nil {
			return 0, 0, err
		}
		limit = int32(limitParse)
	} else {
		limit = 15
	}
	return offset, limit, nil
}
