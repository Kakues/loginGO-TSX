package database

import (
	"fmt"
	"projectGO/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB*gorm.DB

func Connect(){
	
	dsn := "host=localhost user=postgres password=kaiky123 dbname=projectGO port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    v := "Unable to connect to database"
    connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        panic(v)
    }

    DB = connection

    connection.AutoMigrate(&models.User{}, &models.PasswordReset{})
    fmt.Println("Conection ok")
}