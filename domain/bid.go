package domain

type BidNotification struct {
	AuctionID  uint    `json:"auction_id"`
	BidderID   uint    `json:"bidder_id"`
	BidderName string  `json:"bidder_name,omitempty"`
	Value      float64 `json:"value"`
}
