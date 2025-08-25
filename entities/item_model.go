package entities

type Item struct {
	ID             uint             `json:"id" gorm:"primaryKey; autoIncrement"`
	ObjectTypeID   uint             `json:"object_type_id" gorm:"not null;" form:"object_type_id"`
	Name           string           `json:"name" gorm:"not null; unique" form:"name"`
	Price          float64          `json:"price" gorm:"not null" form:"price"`
	DepositPrice   float64          `json:"deposit_price" gorm:"not null" form:"deposit_price"`
	Description    *string          `json:"description" form:"description"`
	FileID         uint             `json:"file_id" gorm:"not null" `
	ObjectType     ObjectType       `json:"object_type,omitempty" gorm:"foreignKey:ObjectTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	File           File             `json:"file,omitempty" gorm:"foreignKey:FileID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ItemDetail     *ItemDetail      `json:"item_detail,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ItemDocument   *ItemDocument    `json:"item_document,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ItemGrade      *ItemGrade       `json:"item_grade,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ItemThumbnails *[]ItemThumbnail `json:"item_thumbnails,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
