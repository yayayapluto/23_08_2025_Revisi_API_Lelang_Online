package entities

type ObjectType struct {
	ID    uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string  `json:"name" gorm:"unique; not null"`
	Items *[]Item `json:"items,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TimeStamp
}
