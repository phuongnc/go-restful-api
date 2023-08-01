package question

import "smartkid/services/common/context"

type UserDomain interface {
	Login(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
}
