package question

import commonDomain "smartkid/services/common/domain/question"

func mapUserFromGorm(userGm *userGorm, _ int) *commonDomain.User {
	if userGm == nil {
		return nil
	}
	return &commonDomain.User{
		Id:        userGm.Id,
		Email:     userGm.Email,
		Name:      userGm.Name,
		Age:       userGm.Age,
		Avatar:    userGm.Avatar,
		Password:  userGm.Password,
		CreatedAt: userGm.CreatedAt,
		UpdatedAt: userGm.UpdatedAt,
	}
}

func mapUserToGorm(user *commonDomain.User) *userGorm {
	if user == nil {
		return nil
	}
	return &userGorm{
		Id:        user.Id,
		Email:     user.Email,
		Password:  user.Password,
		Name:      user.Name,
		Age:       user.Age,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
