package bidderPayment

import (
	"context"
	"errors"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/midtrans"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/user"
)

type (
	Service interface {
		ConfirmedPayment(ctx context.Context, id string) error
		InitializePayment(ctx context.Context, req domain.BidderPaymentRequest) (*string, error)
		FindByBidderData(ctx context.Context, auctionID, userID uint) (*entities.BidderPayment, error)
	}

	service struct {
		repository      Repository
		midtransService midtrans.MidtransService
		userService     user.Service
	}
)

func (s *service) FindByBidderData(ctx context.Context, auctionID, userID uint) (*entities.BidderPayment, error) {
	return s.repository.FindByBidderData(ctx, auctionID, userID)
}

func (s *service) ConfirmedPayment(ctx context.Context, id string) error {
	bidderPayment, err := s.repository.FindById(ctx, id)
	if err != nil {
		return err
	}
	if bidderPayment == nil {
		return errors.New("payment not found")
	}

	bidderPayment.Status = "confirmed"
	if err := s.repository.Update(ctx, bidderPayment); err != nil {
		return errors.New("failed to update payment")
	}

	return nil
}

func (s *service) InitializePayment(ctx context.Context, req domain.BidderPaymentRequest) (*string, error) {
	bidderPayment := entities.BidderPayment{
		BidderID:    uint(req.BidderID),
		Status:      "unknown",
		Amount:      req.Amount,
		RedirectURL: req.RedirectURL,
		Type:        req.Type,
	}

	if err := s.repository.Insert(ctx, &bidderPayment); err != nil {
		return nil, err
	}

	snapUrl, err := s.midtransService.GenerateSnapURL(ctx, &bidderPayment)
	if err != nil {
		return nil, err
	}

	bidderPayment.SnapURL = *snapUrl
	if err := s.repository.Update(ctx, &bidderPayment); err != nil {
		return nil, err
	}

	return &bidderPayment.SnapURL, nil
}

func NewService(repository Repository, midtransService midtrans.MidtransService, userService user.Service) Service {
	return &service{
		repository:      repository,
		midtransService: midtransService,
		userService:     userService,
	}
}
