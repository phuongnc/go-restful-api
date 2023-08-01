package question

import (
	"strings"
	"time"

	commonDomain "smartkid/services/common/domain/question"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRes struct {
	Id                 string         `json:"id"`
	Name               string         `json:"name"`
	Email              string         `json:"email"`
	IsFullVersion      bool           `json:"is_full_version"`
	IsFullVersionAll   bool           `json:"is_full_version_all"`
	NeedUpdatePassword bool           `json:"need_update_password"`
	Avatar             string         `json:"avatar"`
	NeedUpdateInfo     bool           `json:"need_update_info"`
	AccessToken        AccessTokenRes `json:"access_token"`
}

type ShortUserRes struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	TotalPoint       int    `json:"total_point"`
	CurrentChapterId string `json:"current_chapter_id"`
	Avatar           string `json:"avatar"`
}

type AccessTokenRes struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type UpdateUserReq struct {
	Id      string `json:"Id"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Age     int    `json:"age"`
	Country string `json:"country"`
}

func (req LoginReq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Password, validation.Required),
	)
}

func MapUserFromLoginReqDto(loginDto *LoginReq) *commonDomain.User {
	return &commonDomain.User{
		Email:    strings.ToLower(loginDto.Email),
		Password: loginDto.Password,
	}
}

func MapUserToLoginResDto(model *commonDomain.User) *UserRes {
	loginRes := &UserRes{
		Id:                 model.Id,
		Email:              model.Email,
		Name:               model.Name,
		Avatar:             model.Avatar,
		NeedUpdatePassword: model.NeedUpdatePassword,
		NeedUpdateInfo:     false,
	}
	if model.Name == "" {
		loginRes.NeedUpdateInfo = true
	}
	if model.UserSession != nil {
		loginRes.AccessToken = AccessTokenRes{
			Token:     model.UserSession.AccessToken,
			ExpiredAt: model.UserSession.ExpiredAt,
		}
	}
	return loginRes
}

func (req UpdateUserReq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required),
		validation.Field(&req.Avatar, validation.Required),
		validation.Field(&req.Age, validation.Required),
	)
}

func MapUserFromUpdateReqDto(updateUserReq *UpdateUserReq) *commonDomain.User {
	return &commonDomain.User{
		Id:      updateUserReq.Id,
		Name:    updateUserReq.Name,
		Avatar:  updateUserReq.Avatar,
		Age:     updateUserReq.Age,
		Country: updateUserReq.Country,
	}
}

func MapUserToResDto(model *commonDomain.User) *UserRes {
	if model == nil {
		return nil
	}
	userRes := &UserRes{
		Id:                 model.Id,
		Email:              model.Email,
		Name:               model.Name,
		Avatar:             model.Avatar,
		NeedUpdatePassword: model.NeedUpdatePassword,
		NeedUpdateInfo:     false,
	}
	if model.Name == "" {
		userRes.NeedUpdateInfo = true
	}
	return userRes
}
