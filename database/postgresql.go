package database

import (
	"fmt"
	"log"
	"time"

	"github.com/pchawandi/xm-company/config"
	"github.com/pchawandi/xm-company/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	Create(value interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) Database
	Delete(interface{}, ...interface{}) *gorm.DB
	Model(model interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) Database
	Updates(interface{}) *gorm.DB
	Error() error
}

type GormDatabase struct {
	*gorm.DB
}

func (db *GormDatabase) Where(query interface{}, args ...interface{}) Database {
	return &GormDatabase{db.DB.Where(query, args...)}
}

func (db *GormDatabase) First(dest interface{}, conds ...interface{}) Database {
	return &GormDatabase{db.DB.First(dest, conds...)}
}

func (db *GormDatabase) Error() error {
	return db.DB.Error
}

func (db *GormDatabase) Create(value interface{}) *gorm.DB {
	return db.DB.Create(value)
}

func (db *GormDatabase) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return db.DB.Delete(value, conds...)
}

func (db *GormDatabase) Updates(values interface{}) *gorm.DB {
	return db.DB.Updates(values)
}

func (db *GormDatabase) Model(value interface{}) *gorm.DB {
	return db.DB.Model(value)
}

func NewDatabase() *gorm.DB {
	// Load database configuration from the config package
	dbConfig := config.LoadDatabaseConfig()

	// Construct the PostgreSQL connection URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbConfig.User, dbConfig.Password, dbConfig.Hostname, dbConfig.Port, dbConfig.DBName)

	var database *gorm.DB
	var err error

	maxRetries := 3
	for i := 1; i <= maxRetries; i++ {
		database, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err == nil {
			break
		} else {
			log.Printf("Attempt %d: Failed to initialize database. Retrying...", i)
			time.Sleep(3 * time.Second)
		}
	}

	// Ensure the UUID extension is created
	err = database.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Fatal("Failed to create uuid-ossp extension:", err)
	}

	// Ensure the custom ENUM type is created for 'company_type'
	err = database.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'company_type') THEN
				CREATE TYPE company_type AS ENUM ('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship');
			END IF;
		END $$;
	`).Error
	if err != nil {
		log.Print("Failed to create company type enums:", err)
		return nil
	}

	// AutoMigrate models
	_ = database.AutoMigrate(&models.Company{})

	_ = database.AutoMigrate(&models.User{})

	return database
}
