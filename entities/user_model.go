package entities

import (
	"time"
)

type User struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string     `json:"username" gorm:"uniqueIndex;not null"`
	Email       string     `json:"email" gorm:"uniqueIndex;not null"`
	Password    string     `json:"-" gorm:"not null"`
	Role        string     `json:"-" gorm:"type:varchar(5);check:role IN ('user','admin');not null"`
	LastLoginAt *time.Time `json:"last_login_at"`
	TimeStamp
}
