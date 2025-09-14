package entities

import (
	"time"
)

type User struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string     `json:"username" gorm:"uniqueIndex;not null"`
	Email       string     `json:"email" gorm:"uniqueIndex;not null"`
	Password    string     `json:"-" gorm:"not null"`
	Role        string     `json:"-" gorm:"type:varchar(10);check:role IN ('user','organizer','admin');not null"`
	OrganizerID *uint      `json:"organizer_id"`
	Organizer   *Organizer `json:"organizer"`
	LastLoginAt *time.Time `json:"last_login_at"`
	TimeStamp
}
