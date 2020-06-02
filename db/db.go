package db

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/razuos/covid19-bot/config"
	"github.com/razuos/covid19-bot/models"

	// PosgreSQL GORM dialect.
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Connect establishes DB connection.
func Connect() *gorm.DB {
	config.Check([]string{"host", "port", "user", "name", "pass"})
	connectStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", config.AppConfig.DB.Host, config.AppConfig.DB.Port, config.AppConfig.DB.User, config.AppConfig.DB.Name, config.AppConfig.DB.Pass)
	db, err := gorm.Open("postgres", connectStr)
	if err != nil {
		log.Fatalf("Error connecting to DB: %s", err.Error())
	}

	db.AutoMigrate(&models.RSSData{})
	return db
}
