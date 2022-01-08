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
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)
	fmt.Printf("database connection setting: %s\n", dsn)

	globalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
