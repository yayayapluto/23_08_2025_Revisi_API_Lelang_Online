package domain

type (
	BidderPaymentRequest struct {
		Amount      float64 `json:"amount"`
		BidderID    uint64  `json:"bidder_id"`
		RedirectURL string  `json:"redirect_url"`
		Type        string  `json:"type"`
	}
)
