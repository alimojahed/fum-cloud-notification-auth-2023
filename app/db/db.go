package db

import (
	"fmt"
	"log"

	"auth_user/app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) Handler {
	// dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	fmt.Println(url)
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.User{})

	return Handler{db}
}
