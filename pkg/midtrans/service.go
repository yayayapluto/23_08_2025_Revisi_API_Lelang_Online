package midtrans

import (
	"context"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"strconv"
	"strings"
)

type (
	MidtransService interface {
		GenerateSnapURL(ctx context.Context, p *entities.BidderPayment) (*string, error)
		VerifyPayment(ctx context.Context, orderId string) (bool, error)
	}
	midtransService struct {
		config   string
		enviType string
		envi     midtrans.EnvironmentType
	}
)

func (m *midtransService) GenerateSnapURL(ctx context.Context, p *entities.BidderPayment) (*string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(int(p.ID)),
			GrossAmt: int64(p.Amount),
		},
		Callbacks: &snap.Callbacks{
			Finish: "https://instagram.com",
		},
	}

	var client snap.Client
	client.New(m.config, m.envi)

	snapResp, err := client.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	p.SnapURL = snapResp.RedirectURL
	return &p.SnapURL, nil
}

func (m *midtransService) VerifyPayment(ctx context.Context, orderId string) (bool, error) {
	var client coreapi.Client
	client.New(m.config, m.envi)

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := client.CheckTransaction(orderId)
	if e != nil {
		return false, e
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
					return true, nil
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				return true, nil
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				// TODO set transaction status on your databaase to 'failure'
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// TODO set transaction status on your databaase to 'pending' / waiting payment
			}
		}
	}
	return false, nil
}

func NewMidtransService(config, enviType string, envi midtrans.EnvironmentType) MidtransService {
	envi = midtrans.Sandbox
	if strings.ToLower(enviType) == "production" {
		envi = midtrans.Production
	}
	return &midtransService{
		config:   config,
		enviType: enviType,
		envi:     envi,
	}
}
