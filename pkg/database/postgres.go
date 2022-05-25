package database

import (
	"fmt"

	"github.com/qingwave/weave/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(conf *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		conf.Host, conf.User, conf.Password, conf.Name, conf.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
