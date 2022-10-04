package db

import (
	"exampleclean.com/refactor/app/config"
	"exampleclean.com/refactor/app/domain"
	"fmt"
	"github.com/jinzhu/gorm"
)

func ConnectDatabaseGorm(cfg config.Config) (*gorm.DB, error) {
	URL := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName)
	db, err := gorm.Open("mysql", URL)

	db.AutoMigrate(&domain.Users{})

	if err != nil {
		panic(err.Error())
		return nil, err
	}
	return db, nil
}
