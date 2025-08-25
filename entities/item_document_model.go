package entities

type ItemDocument struct {
	ID               uint  `json:"id" gorm:"primaryKey; autoIncrement"`
	ItemID           uint  `json:"item_id" gorm:"not null"`
	Bpkb             *bool `json:"bpkb"`
	Stnk             *bool `json:"stnk"`
	Facture          *bool `json:"facture"`
	Receipt          *bool `json:"receipt"`
	OwnershipRelease *bool `json:"ownership_release"`
	Warranty         *bool `json:"warranty"`
	Box              *bool `json:"box"`

	TimeStamp
}
