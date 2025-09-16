package entities

// AuctionBidder.go
type AuctionBidder struct {
	ID uint `json:"id" gorm:"primaryKey; autoIncrement"`

	UserID        uint   `json:"user_id" gorm:"not null"`
	AuctionID     uint   `json:"auction_id" gorm:"not null"`
	PhoneNumber   string `json:"phone_number" gorm:"not null"`
	BankName      string `json:"bank_name" gorm:"not null"`
	AccountNumber string `json:"account_number" gorm:"not null"`
	AccountName   string `json:"account_name" gorm:"not null"`

	User    User     `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Auction *Auction `json:"auction,omitempty" gorm:"foreignKey:AuctionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	BidderPayment *BidderPayment `json:"bidder_payment,omitempty" gorm:"foreignKey:BidderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Bid           *[]Bid         `json:"bids,omitempty" gorm:"foreignKey:BidderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	TimeStamp
}
