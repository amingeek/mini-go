package main

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	user     = "postgres"
	password = "password"
	host     = "localhost"
	port     = 5432
	dbname   = "postgres"
)

var doOnce sync.Once
var singleton *gorm.DB

func GetConnection() *gorm.DB {
	doOnce.Do(func() {
		connURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			user,
			password,
			host,
			port,
			dbname,
		)

		db, err := gorm.Open(
			postgres.Open(connURL),
			&gorm.Config{
				Logger:                                   logger.Default.LogMode(logger.Silent),
				DisableForeignKeyConstraintWhenMigrating: true,
			},
		)
		if err != nil {
			panic(err)
		}
		singleton = db
	})
	return singleton
}
