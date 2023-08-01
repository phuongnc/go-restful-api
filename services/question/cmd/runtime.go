package cmd

import (
	"fmt"
	"net/http"
	"time"

	commonApp "smartkid/services/common/application/question"
	commonHttp "smartkid/services/common/http"
	"smartkid/services/common/infra/db"
	"smartkid/services/common/infra/logger"
	commonRepo "smartkid/services/common/repository/question"
	"smartkid/services/common/util"
	application "smartkid/services/question/application/http"
	"smartkid/services/question/domain"
	"smartkid/services/question/routing"
)

type runtime struct {
	appConf     *AppConfig
	logger      logger.Logger
	authHandler commonApp.AuthServiceHandler
	userHandler commonApp.UserServiceHandler
	dbc         db.SQL
}

func newRuntime() *runtime {
	rt := runtime{}
	var err error

	if err = util.LoadConfig(&rt.appConf); err != nil {
		panic(fmt.Sprintf("can't load application configuration: %v", err))
	}
	if err = rt.appConf.Validate(); err != nil {
		panic(fmt.Sprintf("application configuration validation failed: %v", err))
	}

	rt.dbc, err = db.NewSQL(rt.appConf.Database)
	if err != nil {
		panic(fmt.Sprintf("cannot connect to db: %v", err))
	}

	rt.logger = logger.NewLogger(rt.appConf.Logger)
	logger.SetDefault(rt.logger)

	baseHandler := commonHttp.NewHandler(rt.logger)
	userRepo := commonRepo.NewUserRepo()
	authDomain := domain.NewAuthDomain(rt.logger, userRepo)
	rt.authHandler = application.NewAuthHandler(baseHandler, rt.logger, authDomain)
	userDomain := domain.NewUserDomain(rt.logger, userRepo)
	rt.userHandler = application.NewUserHandler(baseHandler, rt.logger, userDomain, userRepo)

	return &rt
}

func (rt *runtime) serve() {
	router := routing.NewHTTPRouter(rt.logger, rt.authHandler, rt.dbc, rt.userHandler).SetupHandler()
	readTimeout := time.Minute
	writeTimeout := time.Minute
	endPoint := fmt.Sprintf(":%d", rt.appConf.HTTP.Port)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	rt.logger.Info("[info] start http server listening ", endPoint)
	_ = server.ListenAndServe()
}
