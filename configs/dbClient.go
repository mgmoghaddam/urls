package configs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dbUser := vi.GetString("db.user")
	dbPass := vi.GetString("db.pass")
	dbName := vi.GetString("db.name")
	dbHost := vi.GetString("db.host")
	dbPort := vi.GetString("db.port")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	DB, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
		//panic("failed to connect database")
	}
	DB.LogMode(true)

	//Create db if it doesn't exist
	DB.Exec("CREATE DATABASE IF NOT EXISTS urlshortner")
	if err != nil {
		log.Fatal(err)
	}

	return DB
}

func GetDBClient() *gorm.DB {
	return DB
}
