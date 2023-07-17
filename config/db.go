package config

import (
	"log"
	"sync"

	"github.com/Atgoogat/openmensarobot/db"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	dbConnection     *gorm.DB
	dbConnectionOnce sync.Once
)

func NewDatabaseConnection() (*gorm.DB, error) {
	var e error
	dbConnectionOnce.Do(func() {
		log.Println("connecting to database")
		url, err := Getenv(DB_URL)
		if err != nil {
			e = err
			return
		}

		gormDb, err := gorm.Open(sqlite.Open(url), &gorm.Config{})
		if err != nil {
			e = err
			return
		}
		log.Printf("connected to database: %s", url)
		log.Println("attempting migration")
		err = db.Migrate(gormDb)
		if err != nil {
			e = err
			return
		}
		log.Println("migration complete")
		dbConnection = gormDb
	})

	return dbConnection, e
}
