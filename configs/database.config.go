package configs

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Username           string
	Password           string
	Host               string
	Port               string
	Database           string
	Schema             string
	DatabaseConnection *gorm.DB
}

func (d *DatabaseConfig) connection() {
	dsn := "host=" + d.Host +
		" user=" + d.Username +
		" password=" + d.Password +
		" dbname=" + d.Database +
		" port=" + d.Port +
		" sslmode=disable" +
		" search_path=" + d.Schema // เพิ่ม schema ด้วย search_path

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connect to postgres success!")
	d.DatabaseConnection = db
}

func NewDatabaseConfig() *gorm.DB {
	userName := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")
	schema := os.Getenv("DB_SCHEMA")

	databaseConfig := &DatabaseConfig{
		Username: userName,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
		Schema:   schema,
	}

	databaseConfig.connection()
	return databaseConfig.DatabaseConnection
}
