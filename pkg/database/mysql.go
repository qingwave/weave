package database

import (
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	// Mysql DB
	Mysql *gorm.DB
	once  sync.Once
)

// InitMysql new mysql connection
func InitMysql() error {
	var err error
	once.Do(func() {
		db, dberr := gorm.Open("mysql", "root:123456@tcp(localhost:3306)/weave?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			err = dberr
			return
		}
		Mysql = db
	})
	return err
}
