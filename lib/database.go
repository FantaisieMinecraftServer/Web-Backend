package lib

import (
	"context"
	"log"
	"main/models"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Setup() (*gorm.DB, *mongo.Client, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	dsn := os.Getenv("MYSQL_URI")
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
	err = db.AutoMigrate(
		&models.Account{},
		&models.Player{},
		&models.Economy{},
		&models.Setting{},
	)
	if err != nil {
		return nil, nil, err
	}

	return db, client, nil
}
