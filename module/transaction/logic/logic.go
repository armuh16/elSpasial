package logic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/elspasial/enum"
	"github.com/elspasial/model"
	"github.com/elspasial/module/transaction/dto"
	"github.com/elspasial/module/transaction/repository"
	tripDto "github.com/elspasial/module/trip/dto"
	tripLogic "github.com/elspasial/module/trip/logic"
	userDto "github.com/elspasial/module/user/dto"
	userLogic "github.com/elspasial/module/user/logic"
	"github.com/elspasial/package/logger"
	"github.com/elspasial/static"
	"github.com/elspasial/utilities"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// TransactionLogic
type ITransactionLogic interface {
	CreateOrder(context.Context, *dto.CreateOrderRequest, *gorm.DB) error
	FindAll(context.Context, *dto.FindAllRequest) ([]*model.Transactions, error)
	AcceptOrder(context.Context, *dto.AcceptOrderRequest, *gorm.DB) error
}

type TransactionLogic struct {
	fx.In
	Logger          *logger.LogRus
	TripLogic       tripLogic.IProductLogic
	UserLogic       userLogic.IUserLogic
	TransactionRepo repository.ITransactionRepository
}

// NewLogic :
func NewLogic(transactionLogic TransactionLogic) ITransactionLogic {
	return &transactionLogic
}

func (l *TransactionLogic) CreateOrder(ctx context.Context, reqData *dto.CreateOrderRequest, tx *gorm.DB) error {
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	var (
		grandTotal   float64
		snapshotItem model.TripsTransaction
	)

	userDetail, err := l.UserLogic.Find(ctx, &userDto.FindRequest{
		ID: reqData.UserID,
	})
	if err != nil {
		l.Logger.Error(err)
		return err
	}

	for _, tripID := range reqData.Trip {
		tripDetail, err := l.TripLogic.Find(ctx, &tripDto.FindRequest{
			ID:     tripID,
			UserID: userDetail.ID,
		})
		if err != nil {
			l.Logger.Error(err)
			return err
		}

		snapshotItem = append(snapshotItem, model.TripTransaction{
			ID:          tripDetail.ID,
			Origin:      tripDetail.Origin,
			Destination: tripDetail.Destination,
			Price:       tripDetail.Price,
		})

		grandTotal += tripDetail.Price
	}

	if _, err := l.TransactionRepo.Create(ctx, &model.Transactions{
		UserID:     userDetail.ID,
		GrandTotal: grandTotal,
		Status:     enum.TransactionStatusTypePending,
		Trip:       snapshotItem,
	}, tx); err != nil {
		return utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return nil
}

// FindAll
func (l *TransactionLogic) FindAll(ctx context.Context, reqData *dto.FindAllRequest) ([]*model.Transactions, error) {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	var whereData = model.Transactions{}

	if reqData.RoleID == enum.RoleTypeUser {
		whereData = model.Transactions{
			UserID: reqData.UserID,
		}
	}

	transactions, err := l.TransactionRepo.FindAll(ctx, &whereData)
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "transaksi"), http.StatusNotFound)
		}
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	for _, v := range transactions {
		v.StatusTransaction = v.Status.String()
	}

	return transactions, nil
}

func (l *TransactionLogic) AcceptOrder(ctx context.Context, reqData *dto.AcceptOrderRequest, tx *gorm.DB) error {
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	// Find the pending transaction
	transaction, err := l.TransactionRepo.FindPendingTransaction(ctx, reqData.TransactionID)
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "transaksi"), http.StatusNotFound)
		}
		return utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	// Check if the transaction is already accepted by another driver
	if transaction.Status == enum.TransactionStatusTypeAccept {
		return utilities.ErrorRequest(fmt.Errorf(static.Conflict, "driver"), http.StatusConflict)
	}

	// Update the transaction with the driver ID and status
	transaction.DriverID = reqData.DriverID
	transaction.Status = enum.TransactionStatusTypeAccept

	if err := l.TransactionRepo.Update(ctx, transaction, tx); err != nil {
		l.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return nil
}

// TODO : If status accept and return driver has take the order, if status pending and return driver has not take the order and update the status to accept
// AcceptOrder
//func (l *TransactionLogic) AcceptOrder(ctx context.Context, reqData *dto.AcceptOrderRequest, tx *gorm.DB) error {
//	// Validate request data
//	if err := reqData.Validate(); err != nil {
//		l.Logger.Error(err)
//		return utilities.ErrorRequest(err, http.StatusBadRequest)
//	}
//
//	// Find transaction
//	if _, err := l.TransactionRepo.Find(ctx, &model.Transactions{
//		ID: reqData.TransactionID,
//		//DriverID: reqData.DriverID,
//	}); err != nil {
//		l.Logger.Error(err)
//		if err == gorm.ErrRecordNotFound {
//			return utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "transaksi"), http.StatusNotFound)
//		}
//		return utilities.ErrorRequest(err, http.StatusInternalServerError)
//	}
//
//	if err := l.TransactionRepo.Update(ctx, &model.Transactions{
//		ID:       reqData.TransactionID,
//		DriverID: reqData.DriverID,
//		Status:   enum.TransactionStatusTypeAccept,
//	}, tx); err != nil {
//		l.Logger.Error(err)
//		return utilities.ErrorRequest(err, http.StatusInternalServerError)
//	}
//
//	return nil
//}
