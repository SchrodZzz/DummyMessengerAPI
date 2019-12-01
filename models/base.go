package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"os"
)

var db *gorm.DB

func init() {
	//err := godotenv.Load("mock.env")
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

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&User{}, &Friend{}, &Message{})
}

func GetDB() *gorm.DB {
	return db
}
