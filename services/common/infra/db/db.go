package db

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // blank import
)

// Config config
type Config struct {
	Host               string `config:"host"`
	Port               int    `config:"port"`
	SSL                bool   `config:"ssl"`
	DBMS               string `config:"dbms"`
	Name               string `config:"name"`
	User               string `config:"user"`
	Password           string `config:"password"`
	Schema             string `config:"schema"`
	MaxConnections     int    `config:"max_conn"`
	MaxIdleConnections int    `config:"max_idle_conn"`
	LogMode            bool   `config:"logmode"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
		validation.Field(&c.SSL),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.User),
		validation.Field(&c.Password),
		validation.Field(&c.Schema, validation.Required),
		validation.Field(&c.MaxConnections, validation.Min(1)),
		validation.Field(&c.MaxIdleConnections, validation.Min(1)),
		validation.Field(&c.LogMode),
	)
}

// SQL struct
type SQL interface {
	GetDB() *gorm.DB
	Close() error
}

type sqlStruct struct {
	DB *gorm.DB
}

func (s *sqlStruct) GetDB() *gorm.DB {
	return s.DB
}

func (s *sqlStruct) Close() error {
	if s.DB == nil {
		return nil
	}
	err := s.DB.Close()
	if err != nil {
		return err
	}
	s.DB = nil
	return nil
}

// NewSQL returns new SQL.
func NewSQL(c *Config) (SQL, error) {
	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.Name, c.Schema, func() string {
		if c.SSL {
			return "enable"
		}
		return "disable"
	}())
	db, err := gorm.Open("postgres", connectionStr)

	if err != nil {
		return nil, err
	}

	logrus.Debug("db log mode: ", c.LogMode)
	db.LogMode(c.LogMode)

	if c.MaxConnections == 0 {
		c.MaxConnections = 5
	}
	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = 5
	}

	// pool connection setup
	db.DB().SetMaxOpenConns(c.MaxConnections)
	db.DB().SetMaxIdleConns(c.MaxIdleConnections)
	db.DB().SetConnMaxLifetime(time.Minute * 30 * 8)

	// disable table name's pluralization
	db.SingularTable(true)

	return &sqlStruct{db}, nil
}
