package question

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/samber/lo"

	"smartkid/services/common/context"
	commonDomain "smartkid/services/common/domain/question"
)

type userGorm struct {
	Id                 commonDomain.UUID `gorm:"type:uuid;primary_key"`
	Password           string
	Email              string
	Name               string
	Age                int
	Avatar             string
	IsActive           bool
	IsDeleted          bool
	NeedUpdatePassword bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (d *userGorm) TableName() string {
	return "user"
}

func NewUserRepo() commonDomain.UserRepository {
	return &userRepo{}
}

type userRepo struct{}

func (u *userRepo) Query(ctx context.Context) commonDomain.UserQuery {
	return &userQuery{ctx.GetDBTx().Model(&userGorm{})}
}

func (u *userRepo) InsertUser(ctx context.Context, user *commonDomain.User) error {
	db := ctx.GetDBTx()
	userGm := mapUserToGorm(user)
	if err := db.Create(userGm).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepo) UpdateUser(ctx context.Context, user *commonDomain.User) error {
	db := ctx.GetDBTx()
	userGm := mapUserToGorm(user)
	if err := db.Save(userGm).Error; err != nil {
		return err
	}
	return nil
}

type userQuery struct {
	db *gorm.DB
}

func (u *userQuery) ById(id commonDomain.UUID) commonDomain.UserQuery {
	return &userQuery{db: u.db.Where("id = ?", id)}
}

func (u *userQuery) ByEmail(email string) commonDomain.UserQuery {
	return &userQuery{db: u.db.Where("email = ?", email)}
}

func (u *userQuery) OrderBy(order string) commonDomain.UserQuery {
	return &userQuery{db: u.db.Order(order)}
}

func (u *userQuery) WithPage(offset, limit int32) commonDomain.UserQuery {
	return &userQuery{db: u.db.Offset(offset).Limit(limit)}
}

func (u *userQuery) Result() (*commonDomain.User, error) {
	result := userGorm{}
	err := u.db.Where("is_deleted = ?", false).First(&result).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return mapUserFromGorm(&result, 0), nil
}

func (u *userQuery) ResultList() ([]*commonDomain.User, error) {
	var result []*userGorm
	err := u.db.Where("is_deleted = ?", false).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(result, mapUserFromGorm), nil
}
