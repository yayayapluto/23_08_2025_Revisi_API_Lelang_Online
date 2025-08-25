package entities

type ItemGrade struct {
	ID     uint `json:"id" gorm:"primaryKey; autoIncrement"`
	ItemID uint `json:"item_id" gorm:"not null"`

	Interior string `json:"interior"`
	Exterior string `json:"exterior"`
	Frame    string `json:"frame"`
	Machine  string `json:"machine"`

	TimeStamp
	Item Item `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
