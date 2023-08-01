package question

import "time"

type (
	UUID = string
)

type User struct {
	Id                 UUID
	Name               string
	Email              string
	Password           string
	Age                int
	Country            string
	Avatar             string
	IsActive           bool
	IsDeleted          bool
	NeedUpdatePassword bool
	UserSession        *UserSession
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
