package repository

import (
	"context"

	"github.com/elspasial/database/postgres"
	"github.com/elspasial/model"
	"github.com/elspasial/package/logger"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// TripRepository
type ITripRepository interface {
	Create(context.Context, *model.Trips, *gorm.DB) (*int, error)
	FindAll(context.Context, *model.Trips) ([]*model.Trips, error)
	Find(context.Context, *model.Trips) (*model.Trips, error)
}

type TripRepository struct {
	fx.In
	Logger   *logger.LogRus
	Database *postgres.DB
}

// NewRepository :
func NewRepository(tripRepository TripRepository) ITripRepository {
	return &tripRepository
}

// Create
func (l *TripRepository) Create(ctx context.Context, reqData *model.Trips, tx *gorm.DB) (*int, error) {
	if err := tx.WithContext(ctx).Create(&reqData).Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	return &reqData.ID, nil
}

// FindAll
func (l *TripRepository) FindAll(ctx context.Context, reqData *model.Trips) ([]*model.Trips, error) {
	trips := []*model.Trips{}

	if err := l.Database.Gorm.WithContext(ctx).Model(&model.Trips{}).
		Preload("User").
		Where(&model.Trips{
			UserID: reqData.UserID,
		}).
		Order("id desc").
		Find(&trips).
		Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	return trips, nil
}

// Find
func (l *TripRepository) Find(ctx context.Context, reqData *model.Trips) (*model.Trips, error) {
	product := new(model.Trips)
	if err := l.Database.Gorm.WithContext(ctx).
		Where(&model.Trips{
			ID:     reqData.ID,
			UserID: reqData.UserID,
		}).First(&product).Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	return product, nil
}
