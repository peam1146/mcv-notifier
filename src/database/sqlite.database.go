package database

import (
	"github.com/peam1146/mcv-notifier/src/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSqliteDatabase(dsn string, debug bool) *gorm.DB {

	conf := &gorm.Config{}
	if !debug {
		conf.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(sqlite.Open(dsn), conf)
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&model.Notification{})
	if err != nil {
		panic(err)
	}
	return db
}
