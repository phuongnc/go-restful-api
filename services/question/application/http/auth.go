package http

import (
	"errors"

	"github.com/gin-gonic/gin"

	commonApp "smartkid/services/common/application/question"
	appContext "smartkid/services/common/context"
	commonDomain "smartkid/services/common/domain/question"
	commonHttp "smartkid/services/common/http"
	"smartkid/services/common/infra/logger"
)

type authHandler struct {
	baseHandler          commonHttp.Handler
	logger               logger.Logger
	userDomain           commonDomain.UserDomain
	clientSessionPrivate string
}

func NewAuthHandler(
	baseHandler commonHttp.Handler,
	logger logger.Logger,
	userDomain commonDomain.UserDomain,
) commonApp.AuthServiceHandler {
	return &authHandler{
		baseHandler: baseHandler,
		logger:      logger,
		userDomain:  userDomain,
	}
}

func (h *authHandler) Login(c *gin.Context) {
	clientSessionPublic := c.Request.Header.Get("Client-Session")
	if clientSessionPublic == "" {
		h.logger.Error(errors.New("Client-Session Invalid"))
		h.baseHandler.BadRequest(c, errors.New("Invalid request"))
		return
	}

	loginReq := &commonApp.LoginReq{}
	h.baseHandler.BindToDto(c, loginReq)

	if err := loginReq.Validate(); err != nil {
		h.logger.Error(err, "Validate request fail")
		h.baseHandler.BadRequest(c, err)
		return
	}

	appCtx := appContext.FromContext(c)
	modelUser := commonApp.MapUserFromLoginReqDto(loginReq)
	result, err := h.userDomain.Login(appCtx, modelUser)
	if err != nil {
		h.logger.Error(err, "Login fail")
		h.baseHandler.Unauthorized(c, errors.New("Invalid account"))
		return
	}

	loginRes := commonApp.MapUserToLoginResDto(result)
	h.baseHandler.ResponseData(c, loginRes, "Success")
}
