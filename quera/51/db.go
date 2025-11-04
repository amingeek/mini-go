package main

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	user     = "amin"
	password = "@miN1388"
	host     = "localhost"
	port     = 3306
	dbname   = "quera" // استفاده از دیتابیس quera
)

var doOnce sync.Once
var singleton *gorm.DB

func GetConnection() *gorm.DB {
	doOnce.Do(func() {
		// فرمت DSN برای MySQL
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user,
			password,
			host,
			port,
			dbname,
		)

		db, err := gorm.Open(
			mysql.Open(dsn),
			&gorm.Config{
				Logger:                                   logger.Default.LogMode(logger.Silent),
				DisableForeignKeyConstraintWhenMigrating: false, // تا foreign key ها ساخته بشن
			},
		)

		if err != nil {
			panic(err)
		}
		singleton = db
	})
	return singleton
}
