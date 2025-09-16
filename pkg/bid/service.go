package bid

import (
	"context"
	"encoding/json"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/sse"
	"strconv"
)

type (
	Service interface {
		List(ctx context.Context, offset, limit int, sortDir, sortBy *string) (*[]entities.Bid, int64, error)
		Create(ctx context.Context, e *entities.Bid) (*entities.Bid, error)
		Get(ctx context.Context, id uint) (*entities.Bid, error)
		Update(ctx context.Context, e *entities.Bid) (*entities.Bid, error)
		Delete(ctx context.Context, id uint) error
	}

	service struct {
		repository Repository
		hub        *sse.Hub
	}
)

func (s *service) List(ctx context.Context, offset, limit int, sortDir, sortBy *string) (*[]entities.Bid, int64, error) {
	return s.repository.List(ctx, offset, limit, sortDir, sortBy)
}

func (s *service) Create(ctx context.Context, e *entities.Bid) (*entities.Bid, error) {
	bid, err := s.repository.Create(ctx, e)
	if err != nil {
		return nil, err
	}

	// DTO notif
	notif := domain.BidNotification{
		AuctionID:  bid.Bidder.AuctionID,
		BidderID:   bid.BidderID,
		BidderName: bid.Bidder.User.Username,
		Value:      bid.Value,
	}

	payload, _ := json.Marshal(notif)

	room := "auction_" + strconv.Itoa(int(bid.Bidder.AuctionID))
	s.hub.Broadcast <- sse.RoomMessage{
		Room:    room,
		Message: payload,
	}

	return bid, nil
}

func (s *service) Get(ctx context.Context, id uint) (*entities.Bid, error) {
	return s.repository.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, e *entities.Bid) (*entities.Bid, error) {
	return s.repository.Update(ctx, e)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repository.Delete(ctx, id)
}

func NewService(repository Repository, hub *sse.Hub) Service {
	return &service{repository: repository, hub: hub}
}
