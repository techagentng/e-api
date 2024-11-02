package db

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/techagentng/ecommerce-api/config"
	"github.com/techagentng/ecommerce-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	DB *gorm.DB
}

func GetDB(c *config.Config) *GormDB {
	gormDB := &GormDB{}
	gormDB.Init(c)
	return gormDB
}

func (g *GormDB) Init(c *config.Config) {
	g.DB = getPostgresDB(c)

	if err := migrate(g.DB); err != nil {
		log.Fatalf("unable to run migrations: %v", err)
	}
}

func getPostgresDB(c *config.Config) *gorm.DB {
	log.Printf("Connecting to postgres: %+v", c)
	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=Africa/Lagos",
		c.PostgresHost, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresPort)

	// Create GORM DB instance
	gormConfig := &gorm.Config{}
	if c.Env != "prod" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: postgresDSN,
	}), gormConfig)
	if err != nil {
		log.Fatal(err)
	}

	return gormDB
}

func SeedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{ID: uuid.New(), Name: "Admin"},
		{ID: uuid.New(), Name: "User"},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, models.Role{Name: role.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}

func migrate(db *gorm.DB) error {
		if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
			return fmt.Errorf("failed to create uuid-ossp extension: %v", err)
		}
	// AutoMigrate all the models
	err := db.AutoMigrate(
		&models.User{},
		&models.Blacklist{},
		&models.Role{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		return fmt.Errorf("migrations error: %v", err)
	}

	return nil
}
