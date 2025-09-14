package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"gorm.io/gorm"
	"slices"
	"time"
)

type (
	Repository interface {
		List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.User, int64, error)
		Create(ctx context.Context, user *entities.User) error
		Get(ctx context.Context, id uint) (*entities.User, error)
		GetByEmail(ctx context.Context, email string) (*entities.User, error)
		GetByUsername(ctx context.Context, username string) (*entities.User, error)
		Update(ctx context.Context, user *entities.User) error
		Delete(ctx context.Context, id uint) error

		Register(ctx context.Context, user *entities.User) error
		Login(ctx context.Context, identity, password string) (*string, error)
		GetUserFromToken(ctx context.Context, tokenStr string) (*entities.User, error)
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.User, int64, error) {
	var users []entities.User

	validSortBy := []string{"id", "username", "email", "created_at"}
	defSortBy := validSortBy[0]
	if sortBy != nil {
		if !slices.Contains(validSortBy, *sortBy) {
			return nil, 0, domain.ErrInvalidSortByColumn
		}
		defSortBy = *sortBy
	}

	validSortDir := []string{"asc", "desc"}
	defSortDir := validSortDir[0]
	if sortDir != nil {
		if !slices.Contains(validSortDir, *sortDir) {
			return nil, 0, domain.ErrInvalidSortDir
		}
		defSortDir = *sortDir
	}

	orderStr := fmt.Sprintf("%s %s", defSortBy, defSortDir)
	query := r.db.WithContext(ctx).Model(&entities.User{})

	if search != nil {
		sq := "%" + *search + "%"
		query = query.Where("username LIKE ? OR email LIKE ?", sq, sq)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order(orderStr).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return &users, total, nil
}

func (r *repository) Create(ctx context.Context, user *entities.User) error {
	var count int64

	if user.Username == "" || user.Email == "" {
		return errors.New("username and email are required")
	}

	if err := r.db.WithContext(ctx).Model(&entities.User{}).Where(&entities.User{Username: user.Username, Email: user.Email}).Count(&count).Error; err != nil {
		return err
	}

	if count != 0 {
		return errors.New("username/email already exists")
	}

	newPass, err := utils.HashPassword(user.Email)
	if err != nil {
		return err
	}

	user.Email = newPass
	user.Role = "user"
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) Get(ctx context.Context, id uint) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).First(&user, &entities.User{Email: email}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).First(&user, &entities.User{Username: username}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) Update(ctx context.Context, user *entities.User) error {
	if user.ID == 0 {
		return errors.New("user ID is required for update")
	}

	// Fetch current user to compare old values
	current, err := r.Get(ctx, user.ID)
	if err != nil {
		return err
	}

	updateData := make(map[string]interface{})

	// Check username change
	if user.Username != "" && user.Username != current.Username {
		existingUser, err := r.GetByUsername(ctx, user.Username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if existingUser != nil {
			return gorm.ErrDuplicatedKey
		}
		updateData["username"] = user.Username
	}

	// Check email change
	if user.Email != "" && user.Email != current.Email {
		existingUser, err := r.GetByEmail(ctx, user.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if existingUser != nil {
			return gorm.ErrDuplicatedKey
		}
		updateData["email"] = user.Email
	}

	// If no fields to update, return early
	if len(updateData) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Model(&entities.User{}).
		Where("id = ?", user.ID).
		Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	tx := r.db.WithContext(ctx).Delete(&entities.User{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// AUTH

func (r *repository) Register(ctx context.Context, user *entities.User) error {
	var count int64
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("username, email and password are required")
	}

	// cek unique username/email
	if err := r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("username = ? OR email = ?", user.Username, user.Email).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username/email already exists")
	}

	// hash password
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	user.Role = "user"

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) Login(ctx context.Context, identity, password string) (*string, error) {
	var userModel *entities.User
	var err error

	if utils.CheckIsEmail(identity) {
		userModel, err = r.GetByEmail(ctx, identity)
	} else {
		userModel, err = r.GetByUsername(ctx, identity)
	}

	if err != nil || userModel == nil {
		return nil, errors.New("invalid credentials") // Don't reveal if user exists
	}

	// ✅ Compare provided password with stored hash
	if !utils.CheckPasswordHash(password, userModel.Password) {
		return nil, errors.New("invalid credentials") // Same error for security
	}

	// Generate JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userModel.Username
	claims["role"] = userModel.Role
	claims["user_id"] = userModel.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	env, err := utils.LoadEnv()
	if err != nil {
		return nil, err
	}

	t, err := token.SignedString([]byte(env.JWT_SECRET))
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *repository) GetUserFromToken(ctx context.Context, tokenStr string) (*entities.User, error) {
	env, err := utils.LoadEnv()
	if err != nil {
		return nil, err
	}

	// Parse JWT
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(env.JWT_SECRET), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Ambil user_id
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}
	userID := uint(userIDFloat)

	// Query user dari DB
	user, err := r.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
