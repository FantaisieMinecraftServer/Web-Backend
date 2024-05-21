package lib

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDBConnection() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		panic(err.Error())
	}

	db.SetConnMaxLifetime(20)
	db.SetConnMaxIdleTime(20)

	return db
}

func GetLBDBConnection() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("mysql", os.Getenv("LBDB_URL"))
	if err != nil {
		panic(err.Error())
	}
	return db
}
