package auth

type User struct {
	ID           uint `gorm:"primaryKey"`
	Email        string
	PasswordHash string
}
