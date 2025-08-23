package entities

type ObjectType struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"unique; not null"`
	TimeStamp
}
