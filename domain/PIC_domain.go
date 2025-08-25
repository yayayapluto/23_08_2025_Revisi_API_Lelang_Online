package domain

type UpdateRequestPIC struct {
	Name        *string `json:"name" gorm:"unique; not null"`
	PhoneNumber *string `json:"phone_number" gorm:"unique; not null"`
}
