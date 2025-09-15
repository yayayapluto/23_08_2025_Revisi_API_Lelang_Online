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
	}

	service struct {
		repository      Repository
		midtransService midtrans.MidtransService
		userService     user.Service
	}
)

func (s *service) ConfirmedPayment(ctx context.Context, id string) error {
	bidderPayment, err := s.repository.FindById(ctx, id)
	if err != nil {
		return err
	}

	if bidderPayment == nil {
		return errors.New("payment not found")
	}

	user, err := s.userService.Get(ctx, bidderPayment.BidderID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	return nil
}

func (s *service) InitializePayment(ctx context.Context, req domain.BidderPaymentRequest) (*string, error) {
	bidderPayment := entities.BidderPayment{
		BidderID: uint(req.BidderID),
		Status:   "unknown",
		Amount:   req.Amount,
		Type:     req.Type,
	}

	// insert dulu → auto increment jalan → bidderPayment.ID keisi
	if err := s.repository.Insert(ctx, &bidderPayment); err != nil {
		return nil, err
	}

	// baru generate snap URL pake ID yang udah ada
	if err := s.midtransService.GenerateSnapURL(ctx, &bidderPayment); err != nil {
		return nil, err
	}

	// update snap url ke DB
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
