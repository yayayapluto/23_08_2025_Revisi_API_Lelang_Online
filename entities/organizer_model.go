package entities

type Organizer struct {
	ID            uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string `json:"name" gorm:"unique; not null"`
	Address       string `json:"address" gorm:"not null"`
	BankName      string `json:"bank_name" gorm:"not null"`
	AccountNumber string `json:"account_number" gorm:"not null"`
	AccountName   string `json:"account_name" gorm:"not null"`

	Auctions []Auction `json:"auctions"`
	TimeStamp
}
