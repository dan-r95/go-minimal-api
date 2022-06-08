package login

type User struct {
	ID       uint `gorm:"primaryKey"`
	Email    string
	Password string //hashed
}
