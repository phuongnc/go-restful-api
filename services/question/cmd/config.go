package cmd

import (
	"smartkid/services/common/infra/db"
	"smartkid/services/common/infra/logger"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Client struct {
	Host string `config:"host"`
	Port int    `config:"port"`
	SSL  bool   `config:"ssl"`
}

func (c Client) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Host, validation.Required, is.Host),
		validation.Field(&c.Port, validation.Required, validation.Min(1), validation.Max(65535)),
		validation.Field(&c.SSL),
	)
}

type Storage struct {
	Url string `config:"url"`
}

func (s Storage) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Url, validation.Required),
	)
}

type AppConfig struct {
	HTTP     Client         `config:"http"`
	Database *db.Config     `config:"database"`
	Logger   *logger.Config `config:"logger"`
	Storage  Storage        `config:"storage"`
}

func (c AppConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.HTTP, validation.Required),
		validation.Field(&c.Database, validation.Required),
		validation.Field(&c.Logger, validation.Required),
	)
}
