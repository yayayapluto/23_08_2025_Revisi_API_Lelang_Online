package entities

type File struct {
	ID   uint   `json:"id" gorm:"primaryKey; autoIncrement"`
	Path string `json:"path" gorm:"unique; not null"`
	TimeStamp
}
