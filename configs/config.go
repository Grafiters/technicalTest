package configs

import (
	"log"
	"os"

	"github.com/Grafiters/archive/configs/redises"
	"gorm.io/gorm"
)

var (
	DataBase  *gorm.DB
	Redis     *redises.RedisClient
	JwtConfig *JwtService
	SecretKey string
	Logger    *LoggerFormat
)

func Initialize() error {
	Logger = NewLogger()

	db, err := NewDatabase()
	if err != nil {
		log.Fatal(err)
		return err
	}

	DataBase = db

	Redis, err = redises.NewRedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		return err
	}

	jwt, err := NewJwtConfig()
	if err != nil {
		log.Fatal(err)
		return err
	}
	JwtConfig = jwt

	SecretKey = "skilltest"

	return nil
}
