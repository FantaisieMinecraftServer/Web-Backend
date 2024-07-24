package lib

import (
	"log"
	"main/models"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Setup() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      logger.Warn,
				Colorful:      true,
			},
		),
	})
	if err != nil {
		log.Fatalln(dsn + "database can't connect")
	}

	// Migrate
	err = db.AutoMigrate(&models.Account{}, &models.Player{}, &models.Economy{}, &models.Setting{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
