package domain

import "time"

type UpdateRequestAuction struct {
	ItemID      *uint      `json:"item_id" gorm:"not null"`
	OrganizerID *uint      `json:"organizer_id" gorm:"not null"`
	PicID       *uint      `json:"pic_id" gorm:"not null"`
	StartDate   *time.Time `json:"start_date" gorm:"not null"`
	EndDate     *time.Time `json:"end_date" gorm:"not null"`
}
