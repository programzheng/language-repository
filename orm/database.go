package orm

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	globalDB *gorm.DB
)

func InitDatabase() {
	var err error
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&loc=Local&parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"))

	fmt.Printf("connect: %v database\n", dsn)
	globalDB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256, // default size for string fields
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection Opened to Database")
}

func GetDB() *gorm.DB {
	return globalDB
}

func SetupTableModel(models ...interface{}) error {
	//env is development
	if os.Getenv("APP_ENV") == "development" {
		err := GetDB().AutoMigrate(models...)
		if err != nil {
			log.Fatal(err)
		}
		return err
	}

	return errors.New("")
}
