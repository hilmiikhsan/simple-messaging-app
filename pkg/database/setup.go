package database

import (
	"fmt"
	"log"
	"os"

	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase() (*gorm.DB, *Config) {
	cfg := LoadConfig()
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	DB.Logger = logger.Default.LogMode(logger.Info)

	err = DB.AutoMigrate(&models.User{}, &models.UserSession{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	return DB, cfg
}
