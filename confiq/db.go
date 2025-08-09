package confiq

import (
	"backend/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/go-movie?charset=utf8mb4&parseTime=True&loc=Local"
	database,err:=gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err != nil {
		panic("Gagal koneksi database!")
	}
	fmt.Println("DB connected!")
	database.AutoMigrate(&models.Movie{},&models.User{})
	DB = database
}