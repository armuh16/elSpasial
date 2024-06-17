package logic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/elspasial/model"
	"github.com/elspasial/module/trip/dto"
	"github.com/elspasial/module/trip/repository"
	"github.com/elspasial/package/logger"
	"github.com/elspasial/static"
	"github.com/elspasial/utilities"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProductLogic
type IProductLogic interface {
	Create(context.Context, *dto.CreateRequest, *gorm.DB) error
	FindAll(context.Context, *dto.FindAllRequest) ([]*model.Trips, error)
	Find(context.Context, *dto.FindRequest) (*model.Trips, error)
}

type ProductLogic struct {
	fx.In
	Logger   *logger.LogRus
	TripRepo repository.ITripRepository
}

// NewLogic :
func NewLogic(productLogic ProductLogic) IProductLogic {
	return &productLogic
}

// Create
func (l *ProductLogic) Create(ctx context.Context, reqData *dto.CreateRequest, tx *gorm.DB) error {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	if _, err := l.TripRepo.Create(ctx, &model.Trips{
		UserID:      reqData.UserID,
		Origin:      reqData.Origin,
		Destination: reqData.Destination,
		Price:       reqData.Price,
	}, tx); err != nil {
		return utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return nil
}

// FindAll
func (l *ProductLogic) FindAll(ctx context.Context, reqData *dto.FindAllRequest) ([]*model.Trips, error) {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	trips, err := l.TripRepo.FindAll(ctx, &model.Trips{
		UserID: reqData.UserID,
	})
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "trip"), http.StatusNotFound)
		}
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return trips, nil
}

// FindByID
func (l *ProductLogic) Find(ctx context.Context, reqData *dto.FindRequest) (*model.Trips, error) {
	product, err := l.TripRepo.Find(ctx, &model.Trips{
		ID:     reqData.ID,
		UserID: reqData.UserID,
	})
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "trip"), http.StatusNotFound)
		}
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}
	return product, nil
}
