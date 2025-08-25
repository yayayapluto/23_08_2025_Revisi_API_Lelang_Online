package entities

type ItemThumbnail struct {
	ID     uint   `json:"id" gorm:"primaryKey; autoIncrement"`
	Name   string `json:"name" gorm:"not null" form:"name"`
	ItemID uint   `json:"item_id" gorm:"not null" form:"item_id"`
	FileID uint   `json:"file_id" gorm:"not null"`

	TimeStamp
	Item Item `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	File File `gorm:"foreignKey:FileID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
