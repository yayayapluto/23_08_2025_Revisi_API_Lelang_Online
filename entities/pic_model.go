package entities

type PIC struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"unique; not null"`
	PhoneNumber string `json:"phone_number" gorm:"unique; not null"`

	Auctions *[]Auction `json:"auctions" gorm:"foreignKey:PicID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TimeStamp
}
