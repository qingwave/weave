package database

import (
	"fmt"
	"weave/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(conf *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		conf.DBHost, conf.User, conf.Password, conf.DBName, conf.DBPort)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
