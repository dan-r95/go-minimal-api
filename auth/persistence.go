package auth

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"regexp"
	"strconv"
)

var (
	ErrCredentials   = errors.New("credentials dont match or email empty")
	RExpPass         = regexp.MustCompile("^.{6,}")
	ErrPasswordMatch = errors.New("wrong credentials")
	RExpMail         = regexp.MustCompile("^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)+$")
)

// database defines the persistence layer, including creation and look up of users
type database struct {
	db *gorm.DB
	*zap.SugaredLogger
}

func (d *database) setup() (db *gorm.DB, err error) {
	port, err := strconv.ParseInt(os.Getenv("POSTGRES_PORT"), 10, 64)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"), port,
	)
	d.Infow("db env", "env", dsn)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}
	// Migrate the schema
	if err = db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	if err = db.Create(&User{PasswordHash: "D42", Email: "test-mail@mtest.de"}).Error; err != nil {
		return nil, err
	}
	return
}

func (d *database) Login(req LoginRequest) error {
	user := &User{}
	if req.Email == "" || req.Password == "" {
		return ErrCredentials
	}
	if err := d.db.Take(user, "email = ?", req.Email).Error; err != nil {
		return err
	} else if err = verify(user.PasswordHash, req.Password); err != nil {
		return ErrPasswordMatch
	} else {
		// success
		return nil
	}
}

func (d *database) Register(req RegistrationRequest) error {
	if req.Email == "" || !RExpMail.MatchString(req.Email) || req.Password != req.Confirmation {
		return ErrCredentials
	} else if hash, err := createHash(req.Password); err != nil {
		return err
	} else {
		return d.db.Transaction(func(tx *gorm.DB) error {
			user := &User{Email: req.Email, PasswordHash: hash}
			return d.db.Create(user).Error
		})
	}
}

func createHash(password string) (string, error) {
	if password == "" || !RExpPass.MatchString(password) {
		return "", ErrCredentials
	}

	if b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		// a password which cannot be encrypted has to lead to a panic
		panic(fmt.Errorf("encryption of password failed: %s", err.Error()))
	} else {
		return string(b), nil
	}
}

func verify(passwordHash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
