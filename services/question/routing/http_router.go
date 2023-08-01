package routing

import (
	"smartkid/services/common/infra/db"
	"smartkid/services/common/infra/logger"
	"smartkid/services/common/middleware"

	commonApp "smartkid/services/common/application/question"

	"github.com/gin-gonic/gin"
)

type HTTPRouter interface {
	SetupHandler() *gin.Engine
}

type httpRouter struct {
	logger      logger.Logger
	authHandler commonApp.AuthServiceHandler
	userHandler commonApp.UserServiceHandler
	db          db.SQL
}

func (hr *httpRouter) SetupHandler() *gin.Engine {
	var router = gin.Default()
	router.Use(middleware.Cors())
	router.Use(middleware.HttpDb(hr.db, hr.logger))

	ApiGroup := router.Group("/v1")
	authRouter := ApiGroup.Group("/auth")
	{
		authRouter.POST("login", hr.authHandler.Login)
	}
	userRouter := ApiGroup.Group("/users").Use(middleware.JWTAuth(hr.logger))
	{
		userRouter.GET("/user-info", hr.userHandler.GetUser)
		userRouter.PUT("/:user_id", hr.userHandler.UpdateUser)
	}
	return router
}

func NewHTTPRouter(
	logger logger.Logger,
	authHandler commonApp.AuthServiceHandler,
	db db.SQL,
	userHandler commonApp.UserServiceHandler,
) HTTPRouter {
	return &httpRouter{
		logger:      logger,
		authHandler: authHandler,
		userHandler: userHandler,
		db:          db,
	}
}
