package question

import "time"

type UserSession struct {
	Id          UUID
	UserId      UUID
	AccessToken string
	ExpiredAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
