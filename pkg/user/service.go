package user

import (
	"context"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service interface {
		List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.User, int64, error)
		Create(ctx context.Context, user *entities.User) error
		Get(ctx context.Context, id uint) (*entities.User, error)
		GetByEmail(ctx context.Context, email string) (*entities.User, error)
		GetByUsername(ctx context.Context, username string) (*entities.User, error)
		Update(ctx context.Context, user *entities.User) error
		Delete(ctx context.Context, id uint) error

		CreateOrganizerUser(ctx context.Context, input *domain.CreateOrganizerUserInput) error

		Register(ctx context.Context, input *domain.RegisterInput) error
		Login(ctx context.Context, identity, password string) (*string, error)
		GetUserFromToken(ctx context.Context, tokenStr string) (*entities.User, error)
	}

	service struct {
		repo Repository
	}
)

func (s *service) List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.User, int64, error) {
	return s.repo.List(ctx, offset, limit, search, sortDir, sortBy)
}

func (s *service) Create(ctx context.Context, user *entities.User) error {
	return s.repo.Create(ctx, user)
}

func (s *service) Get(ctx context.Context, id uint) (*entities.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *service) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *service) Update(ctx context.Context, user *entities.User) error {
	return s.repo.Update(ctx, user)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) CreateOrganizerUser(ctx context.Context, input *domain.CreateOrganizerUserInput) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := entities.User{
		Username:    input.Username,
		Email:       input.Email,
		Password:    string(hashedPwd),
		Role:        "organizer",
		OrganizerID: &input.OrganizerID,
	}

	return s.repo.Create(ctx, &user)
}

// auth
func (s *service) Register(ctx context.Context, input *domain.RegisterInput) error {
	user := &entities.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password, // hash di repo
		Role:     "user",
	}
	return s.repo.Register(ctx, user)
}

func (s *service) Login(ctx context.Context, identity, password string) (*string, error) {
	return s.repo.Login(ctx, identity, password)
}

func (s *service) GetUserFromToken(ctx context.Context, tokenStr string) (*entities.User, error) {
	return s.repo.GetUserFromToken(ctx, tokenStr)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
