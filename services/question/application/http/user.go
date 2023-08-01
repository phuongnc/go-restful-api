package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	commonApp "smartkid/services/common/application/question"
	appContext "smartkid/services/common/context"
	commonDomain "smartkid/services/common/domain/question"
	commonHttp "smartkid/services/common/http"
	"smartkid/services/common/infra/logger"
)

type userHandler struct {
	baseHandler commonHttp.Handler
	logger      logger.Logger
	userDomain  commonDomain.UserDomain
	userRepo    commonDomain.UserRepository
}

func NewUserHandler(
	baseHandler commonHttp.Handler,
	logger logger.Logger,
	userDomain commonDomain.UserDomain,
	userRepo commonDomain.UserRepository,
) commonApp.UserServiceHandler {
	return &userHandler{
		baseHandler: baseHandler,
		logger:      logger,
		userDomain:  userDomain,
		userRepo:    userRepo,
	}
}

func (h *userHandler) UpdateUser(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		h.logger.Error("User ID is empty")
		h.baseHandler.BadRequest(c, nil, "Invalid param")
		return
	}

	updateUserReq := &commonApp.UpdateUserReq{}
	h.baseHandler.BindToDto(c, updateUserReq)
	updateUserReq.Id = userId

	if err := updateUserReq.Validate(); err != nil {
		h.logger.Error(err, "Validate request fail")
		h.baseHandler.BadRequest(c, err)
		return
	}

	appCtx := appContext.FromContext(c)
	modelUser := commonApp.MapUserFromUpdateReqDto(updateUserReq)
	result, err := h.userDomain.UpdateUser(appCtx, modelUser)
	if err != nil {
		h.logger.Error(err, "Can not update user")
		h.baseHandler.InternalServerError(c, errors.New("Can not update user info"))
		return
	}

	updateInfoRes := commonApp.MapUserToResDto(result)
	h.baseHandler.ResponseData(c, updateInfoRes, "Success")
}

func (h *userHandler) GetUser(c *gin.Context) {
	appCtx := appContext.FromContext(c)
	userId := appCtx.GetUserId()
	result, err := h.userRepo.Query(appCtx).ById(userId).Result()
	if err != nil {
		h.logger.Error(err, "Can not get user")
		h.baseHandler.InternalServerError(c, errors.New("Can not get user info"))
		return
	}
	if result == nil {
		h.logger.Error(err, "User is not exist")
		h.baseHandler.ResponseError(c, "User not found", nil, http.StatusNotFound)
		return
	}

	updateInfoRes := commonApp.MapUserToResDto(result)
	h.baseHandler.ResponseData(c, updateInfoRes, "Success")
}
