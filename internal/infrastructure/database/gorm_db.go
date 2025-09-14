package database

import (
	"fmt"
	"log"
	"time"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/model"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormConnection(connStr string, debug bool) *gorm.DB {
	gormConfig := &gorm.Config{}
	if debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	var db *gorm.DB
	var err error

	for {
		fmt.Println("Connecting to SQL-DB...")
		var dbConf = sqlite.Open("debug.db")
		if !debug {
			dbConf = postgres.Open(connStr)
		}
		db, err = gorm.Open(dbConf, gormConfig)
		if err != nil {
			fmt.Println("failed to connect database:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	err = db.AutoMigrate(
		&model.Log{},
	)
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	return db
}
