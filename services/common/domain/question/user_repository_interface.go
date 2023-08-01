package question

import "smartkid/services/common/context"

type UserRepository interface {
	Query(ctx context.Context) UserQuery
	InsertUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
}

type UserQuery interface {
	ById(id UUID) UserQuery
	ByEmail(email string) UserQuery
	OrderBy(order string) UserQuery
	WithPage(offset, limit int32) UserQuery
	ResultList() ([]*User, error)
	Result() (*User, error)
}
