package repository

import (
	"context"
	"gorm.io/gorm"

	"github.com/elspasial/database/postgres"
	"github.com/elspasial/model"
	"github.com/elspasial/package/logger"

	"go.uber.org/fx"
)

// UserRepository
type IUserRepository interface {
	Find(context.Context, *model.Users) (*model.Users, error)
	Create(context.Context, *model.Users, *gorm.DB) error
	CreateReturnId(ctx context.Context, users *model.Users, db *gorm.DB) (*int, error)
}

type UserRepository struct {
	fx.In
	Logger   *logger.LogRus
	Database *postgres.DB
}

// NewRepository :
func NewRepository(userRepository UserRepository) IUserRepository {
	return &userRepository
}

// create and return id
func (l *UserRepository) CreateReturnId(ctx context.Context, users *model.Users, db *gorm.DB) (*int, error) {
	err := l.Create(ctx, users, db)
	if err != nil {
		return nil, err
	}

	return &users.ID, nil
}

// Create
func (l *UserRepository) Create(ctx context.Context, reqData *model.Users, tx *gorm.DB) error {
	if err := tx.WithContext(ctx).Create(&reqData).Error; err != nil {
		l.Logger.Error(err)
		return err
	}
	return nil
}

// Find
func (l *UserRepository) Find(ctx context.Context, reqData *model.Users) (*model.Users, error) {
	product := new(model.Users)

	if err := l.Database.Gorm.WithContext(ctx).
		Where(&model.Users{
			ID:    reqData.ID,
			Email: reqData.Email,
		}).First(&product).Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	return product, nil
}
