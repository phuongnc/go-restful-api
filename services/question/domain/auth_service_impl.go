package domain

import (
	"errors"
	"time"

	"smartkid/services/common/context"
	domain "smartkid/services/common/domain/question"
	"smartkid/services/common/infra/logger"
	"smartkid/services/common/util"

	"github.com/google/uuid"
)

type authDomainImpl struct {
	logger   logger.Logger
	userRepo domain.UserRepository
}

func NewAuthDomain(logger logger.Logger, userRepo domain.UserRepository) domain.UserDomain {
	return &authDomainImpl{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (a *authDomainImpl) Login(ctx context.Context, user *domain.User) (*domain.User, error) {
	a.logger.DebugCtx(ctx, "Get user by email")
	result, err := a.userRepo.Query(ctx).ByEmail(user.Email).Result()
	if err != nil {
		a.logger.ErrorCtx(ctx, err, "Get user by email fail")
		return nil, err
	}
	if result == nil {
		//create use here
		passwordHash, err := util.HashPassword(user.Password)
		if err != nil {
			a.logger.ErrorCtx(ctx, err, "Error when hash password")
			return nil, errors.New("Error when hash password")
		}
		user.Id = uuid.New().String()
		user.Password = passwordHash
		user.CreatedAt = time.Now()

		err = a.userRepo.InsertUser(ctx, user)
		if err != nil {
			a.logger.ErrorCtx(ctx, err, "Can not create new User")
			return nil, errors.New("Can not create new User")
		}
		result, err = a.userRepo.Query(ctx).ById(user.Id).Result()
		if err != nil {
			a.logger.ErrorCtx(ctx, err, "Can not get new User")
			return nil, errors.New("Can not get new User")
		}
	} else {
		if !util.CheckPasswordHash(user.Password, result.Password) {
			a.logger.ErrorCtx(ctx, err, "Email or Password wrong")
			return nil, errors.New("Login information is wrong, please try again")
		}
	}

	j := util.NewJWT()
	tokenInfo, err := j.GenerateToken(result.Id, result.Email)
	if err != nil {
		a.logger.ErrorCtx(ctx, err, "Can not generate token")
		return nil, errors.New("Can not generate token")
	}

	userSession := &domain.UserSession{
		Id:          uuid.New().String(),
		UserId:      result.Id,
		AccessToken: tokenInfo.Token,
		ExpiredAt:   tokenInfo.ExpiredAt,
		CreatedAt:   time.Now(),
	}

	result.UserSession = userSession
	return result, nil
}

func (a *authDomainImpl) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	panic("No need implement")
}
