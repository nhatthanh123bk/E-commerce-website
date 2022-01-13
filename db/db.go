package db

import (
	"github.com/blogs/app/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/blogs/helper"

)

var DB *gorm.DB
var err error

func Init() {
	// dbUrl := os.Getenv("DB_URL")
	dbUrl := "nhatthanh:Baochau@2809@tcp(192.168.56.10:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"

	if DB, err = gorm.Open(mysql.Open(dbUrl)); err != nil {
		panic(err)
	}
	helper.Logger.Infow("Connected to database susscessfully!")

	// auto migrate all of available models
	initMigrate()
}

func initMigrate() {
	DB.AutoMigrate(&model.User{})
}
