package logic

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"

	"github.com/elspasial/model"
	"github.com/elspasial/module/user/dto"
	"github.com/elspasial/module/user/repository"
	"github.com/elspasial/package/logger"
	"github.com/elspasial/static"
	"github.com/elspasial/utilities"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// UserLogic
type IUserLogic interface {
	Find(context.Context, *dto.FindRequest) (*model.Users, error)
	Create(context.Context, *dto.CreateRequest, *gorm.DB) (*int, error)
}

type UserLogic struct {
	fx.In
	Logger   *logger.LogRus
	UserRepo repository.IUserRepository
}

// NewLogic :
func NewLogic(userLogic UserLogic) IUserLogic {
	return &userLogic
}

// Create
func (l *UserLogic) Create(ctx context.Context, reqData *dto.CreateRequest, tx *gorm.DB) (*int, error) {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	// Check if email has taken
	existEmail, err := l.UserRepo.Find(ctx, &model.Users{
		Email: reqData.Email,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	} else if existEmail != nil {
		return nil, utilities.ErrorRequest(fmt.Errorf(static.ValueAreadyExist, "Email"), http.StatusBadRequest)
	}

	// Hash password
	password, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	// Create new user
	reqData.Password = string(password)
	userId, err := l.UserRepo.CreateReturnId(ctx, &model.Users{
		Name:     reqData.Name,
		Password: string(password),
		Email:    strings.ToLower(reqData.Email),
		Address:  reqData.Adress,
		RoleID:   reqData.RoleID,
	}, tx)
	if err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return userId, nil
}

// FindByID
func (l *UserLogic) Find(ctx context.Context, reqData *dto.FindRequest) (*model.Users, error) {
	product, err := l.UserRepo.Find(ctx, &model.Users{
		ID:    reqData.ID,
		Email: reqData.Email,
	})
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "user"), http.StatusNotFound)
		}
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}
	return product, nil
}
