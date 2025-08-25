package entities

type ItemThumbnail struct {
	ID     uint   `json:"id" gorm:"primaryKey; autoIncrement"`
	Name   string `json:"name" gorm:"not null" form:"name"`
	ItemID uint   `json:"item_id" gorm:"not null" form:"item_id"`
	FileID uint   `json:"file_id" gorm:"not null"`

	TimeStamp
	File File `json:"file,omitempty" gorm:"foreignKey:FileID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
