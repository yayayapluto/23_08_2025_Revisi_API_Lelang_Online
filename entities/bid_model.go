package entities

type Bid struct {
	ID uint `json:"id" gorm:"primaryKey; autoIncrement"`

	BidderID uint    `json:"bidder_id" gorm:"not null"`
	Value    float64 `json:"value" gorm:"not null"`

	Bidder AuctionBidder `json:"bidder,omitempty" gorm:"foreignKey:BidderID"`

	TimeStamp
}
