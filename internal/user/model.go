package user

import "time"

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
}
