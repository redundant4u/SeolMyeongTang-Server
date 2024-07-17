package models

import (
	"fmt"

	"github.com/redundant4u/SeolMyeongTang-Server/chat/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable timezone=Asia/Seoul", configs.Env.DBHost, configs.Env.DBUser, configs.Env.DBPassword, configs.Env.DBNAME, configs.Env.DBPort)
	pg, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("DB connection failed")
	}

	DB = pg
}
