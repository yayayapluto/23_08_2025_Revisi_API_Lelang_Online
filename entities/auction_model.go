package entities

import "time"

type Auction struct {
	ID          uint      `json:"id" gorm:"primaryKey; autoIncrement"`
	ItemID      uint      `json:"item_id" gorm:"not null"`
	OrganizerID uint      `json:"organizer_id" gorm:"not null"`
	PicID       uint      `json:"pic_id" gorm:"not null"`
	StartDate   time.Time `json:"start_date" gorm:"not null"`
	EndDate     time.Time `json:"end_date" gorm:"not null"`
	TimeStamp

	Item      Item      `json:"item" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Organizer Organizer `json:"organizer" gorm:"foreignKey:OrganizerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PIC       PIC       `json:"pic" gorm:"foreignKey:PicID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Bidders []AuctionBidder `json:"bidders" gorm:"foreignKey:AuctionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
