package repository

import (
	"context"
	"github.com/elspasial/database/postgres"
	"github.com/elspasial/enum"
	"github.com/elspasial/model"
	"github.com/elspasial/package/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// TransactionRepository
type ITransactionRepository interface {
	Create(context.Context, *model.Transactions, *gorm.DB) (*int, error)
	FindAll(context.Context, *model.Transactions) ([]*model.Transactions, error)
	Find(context.Context, *model.Transactions) (*model.Transactions, error)
	Update(context.Context, *model.Transactions, *gorm.DB) error
	FindPendingTransaction(context.Context, int) (*model.Transactions, error)
}

type TransactionRepository struct {
	fx.In
	Logger   *logger.LogRus
	Database *postgres.DB
}

// NewRepository :
func NewRepository(transactionRepository TransactionRepository) ITransactionRepository {
	return &transactionRepository
}

// Create
func (l *TransactionRepository) Create(ctx context.Context, reqData *model.Transactions, tx *gorm.DB) (*int, error) {
	if err := tx.WithContext(ctx).Create(&reqData).Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	return &reqData.ID, nil
}

// FindAll
func (l *TransactionRepository) FindAll(ctx context.Context, reqData *model.Transactions) ([]*model.Transactions, error) {
	transactions := []*model.Transactions{}

	if err := l.Database.Gorm.WithContext(ctx).Model(&model.Transactions{}).
		Preload("User").
		Where(&model.Transactions{
			UserID: reqData.UserID,
		}).
		Order("id desc").
		Find(&transactions).
		Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	return transactions, nil
}

// Find
func (l *TransactionRepository) Find(ctx context.Context, reqData *model.Transactions) (*model.Transactions, error) {
	transaction := new(model.Transactions)
	if err := l.Database.Gorm.WithContext(ctx).
		Where(&model.Transactions{
			ID: reqData.ID,
			//DriverID: reqData.DriverID,
		}).
		First(&transaction).Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	return transaction, nil
}

func (r *TransactionRepository) FindPendingTransaction(ctx context.Context, transactionID int) (*model.Transactions, error) {
	transaction := new(model.Transactions)
	if err := r.Database.Gorm.WithContext(ctx).
		Where("id = ? AND driver_id = 0 AND status = ?", transactionID, enum.TransactionStatusTypePending).
		First(transaction).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}
	return transaction, nil
}

func (r *TransactionRepository) Update(ctx context.Context, reqData *model.Transactions, tx *gorm.DB) error {
	if err := tx.WithContext(ctx).Model(&model.Transactions{}).
		Where("id = ?", reqData.ID).
		Updates(reqData).Error; err != nil {
		r.Logger.Error(err)
		return err
	}
	return nil
}

// Update
//func (l *TransactionRepository) Update(ctx context.Context, reqData *model.Transactions, tx *gorm.DB) error {
//	if err := tx.WithContext(ctx).Model(&model.Transactions{}).
//		Where("id = ?", reqData.ID).
//		Where("driver_id = ?", reqData.DriverID).
//		Updates(model.Transactions{
//			DriverID:  reqData.DriverID,
//			Status:    reqData.Status,
//			UpdatedAt: time.Now(),
//		}).Error; err != nil {
//		l.Logger.Error(err)
//		return err
//	}
//	return nil
//}
