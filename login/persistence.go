package login

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// database defines the persistence layer, including creation and look up of users
type database struct {
	db *gorm.DB
}

func (d database) create() (err error) {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	d.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.New("failed to connect database")
	}
	return
}
