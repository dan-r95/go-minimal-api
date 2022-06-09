package auth

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
)

// database defines the persistence layer, including creation and look up of users
type database struct {
	db *gorm.DB
	*zap.SugaredLogger
}

func (d *database) create() (err error) {
	port, err := strconv.ParseInt(os.Getenv("POSTGRES_PORT"), 10, 64)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"), port,
	)
	d.Infow("db env", "env", dsn)
	d.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.New("failed to connect database")
	}
	return
}
