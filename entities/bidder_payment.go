package entities

// BidderPayment.go
type BidderPayment struct {
	ID uint `json:"id" gorm:"primaryKey; autoIncrement"`

	BidderID uint    `json:"bidder_id" gorm:"not null;unique"` // unique biar one-to-one
	Status   string  `json:"status" gorm:"not null"`
	Amount   float64 `json:"amount" gorm:"not null"`
	SnapURL  string  `json:"snap_url" gorm:"not null"`
	Type     string  `json:"type" gorm:"not null"` // deposit or final

	Bidder AuctionBidder `json:"bidder,omitempty" gorm:"foreignKey:BidderID"`

	TimeStamp
}
