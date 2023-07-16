package config

import (
	"log"
	"os"

	"github.com/Atgoogat/openmensarobot/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabaseConnection() *gorm.DB {
	url, ok := os.LookupEnv(DB_URL)
	if !ok {
		log.Fatalf("Env var %s is not defined", DB_URL)
	}

	gormDb, err := gorm.Open(sqlite.Open(url), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to database: %s\n", url)
	log.Println("Attempting migration")
	err = db.Migrate(gormDb)
	if err != nil {
		panic(err)
	}
	log.Println("Migration complete")

	return gormDb
}
