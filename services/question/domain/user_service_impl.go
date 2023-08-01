package domain

import (
	"errors"
	"time"

	"smartkid/services/common/context"
	domain "smartkid/services/common/domain/question"
	"smartkid/services/common/infra/logger"
)

type userDomainImpl struct {
	logger   logger.Logger
	userRepo domain.UserRepository
}

func NewUserDomain(logger logger.Logger, userRepo domain.UserRepository) domain.UserDomain {
	return &userDomainImpl{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (a *userDomainImpl) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	a.logger.DebugCtx(ctx, "Checking the user")
	existingUser, err := a.userRepo.Query(ctx).ById(user.Id).Result()
	if err != nil {
		a.logger.ErrorCtx(ctx, err, "Can not get User by id=", user.Id)
		return nil, err
	}
	if existingUser == nil {
		a.logger.ErrorCtx(ctx, err, "User is not exist")
		return nil, errors.New("Invalid User")
	}
	a.logger.DebugCtx(ctx, "Update user info")

	existingUser.Age = user.Age
	existingUser.Avatar = user.Avatar
	existingUser.Name = user.Name
	existingUser.UpdatedAt = time.Now()
	existingUser.Country = user.Country
	err = a.userRepo.UpdateUser(ctx, existingUser)
	if err != nil {
		a.logger.ErrorCtx(ctx, err, "Can not update user")
		return nil, err
	}

	return existingUser, nil
}

func (a *userDomainImpl) Login(ctx context.Context, user *domain.User) (*domain.User, error) {
	panic("No need to implement")
}
