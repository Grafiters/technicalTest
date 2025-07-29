package configs

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	var dialector gorm.Dialector

	var sslmode string
	if os.Getenv("DB_SSLMODE") == "disable" {
		sslmode = "disable"
	} else {
		sslmode = "require"
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASS") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + sslmode

	dialector = postgres.Open(dsn)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
