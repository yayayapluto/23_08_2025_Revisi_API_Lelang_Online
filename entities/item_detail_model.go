package entities

import "time"

type ItemDetail struct {
	ID            uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ItemID        uint       `json:"item_id" gorm:"not null"`
	PlateNumber   *string    `json:"plate_number"`
	Brand         string     `json:"brand" gorm:"not null"`
	Series        *string    `json:"series"`
	CC            *float64   `json:"cc"`
	Type          *string    `json:"type"`
	Transmission  *string    `json:"transmission"`
	Model         *string    `json:"model"`
	Year          int        `json:"year" gorm:"not null"`
	FrameNumber   *string    `json:"frame_number"`
	MachineNumber *string    `json:"machine_number"`
	Kilometer     *int       `json:"kilometer"`
	Fuel          *string    `json:"fuel"`
	Color         string     `json:"color" gorm:"not null"`
	DriveType     *string    `json:"drive_type"`
	StnkDate      *time.Time `json:"stnk_date"`

	TimeStamp

	Item Item `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
