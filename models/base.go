package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
)

var db *gorm.DB

func init() {
	//err := godotenv.Load(".env")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//name := os.Getenv("db_name")
	//pass := os.Getenv("db_pass")
	//port := os.Getenv("db_port")
	//schema := os.Getenv("db_schema")
	//
	//path := fmt.Sprintf("%s:%s@(localhost:%s)/%s?charset=utf8&parseTime=True&loc=Local", name, pass, port, schema)
	//
	//db, err = gorm.Open("mysql", path)
	//if err != nil {
	//	fmt.Println(err)
	//}

	err := godotenv.Load() //Load .env file
	if err != nil {
		fmt.Print(err)
	}

	db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Print(err)
	}

	db.Debug().AutoMigrate(&User{}, &Friend{}, &Message{})
}

func GetDB() *gorm.DB {
	return db
}
