package database

import (
	"fmt"
	"log"
	"m/v2/app/models"
	"os"
	"time"

	//"github.com/RCAFAWC/maintainer/app/models"
	//"github.com/RCAFAWC/maintainer/config"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	// DBConn for gorm
	DBConn *gorm.DB
	// SessionStore for session storage
	SessionStore *session.Store
	cfg          *config.DatabaseConfig
)

type Database struct {
	*gorm.DB
}

// InitGormDb with gorm models
func Setup() {
	if DBConn != nil {
		return
	}

	cfg = config.GetInstance().GetDatabaseConfig()

	if err := connect(); err != nil {
		log.Panicf("error could not connect database (%s)", err.Error())
	}

	log.Println("Connected to DB")

	if err := DBConn.AutoMigrate(&models.Task{},
		&models.User{}); err != nil {
		log.Panicf("Failed to automigrate %s", err.Error())
	}

	log.Println("Ran Auto Migrate")
}

func connect() error {
	connectionString := ""
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	var err error
	switch cfg.Default.Driver {
	case "postgresql", "postgres":
		connectionString = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", cfg.Default.Host, cfg.Default.Port, cfg.Default.Username, cfg.Default.DBName, cfg.Default.Password)
		DBConn, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
		})
	default:
		log.Fatalf("error: database driver '%s' not supported", cfg.Default.Driver)
	}
	if err != nil {
		fmt.Println(cfg.Default)
		panic(err)
	}
	if err = DBConn.Use(
		dbresolver.Register(dbresolver.Config{}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	); err != nil {
		log.Fatalf("error: registering resolver failed '%s'", err.Error())
	}
	return nil
}
