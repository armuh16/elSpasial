package route

import (
	"errors"
	"net/http"

	"github.com/elspasial/database/postgres"
	"github.com/elspasial/module/transaction/dto"
	"github.com/elspasial/module/transaction/logic"
	"github.com/elspasial/package/jwt"
	"github.com/elspasial/package/logger"
	"github.com/elspasial/router"
	"github.com/elspasial/static"
	"github.com/elspasial/utilities"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Handler defines the handler for transaction routes
type Handler struct {
	fx.In
	Logic     logic.ITransactionLogic
	EchoRoute *router.Router
	Logger    *logger.LogRus
	Db        *postgres.DB
}

func NewRoute(h Handler, m ...echo.MiddlewareFunc) Handler {
	h.Route(m...)
	return h
}

func (h *Handler) Route(m ...echo.MiddlewareFunc) {
	transaction := h.EchoRoute.Group("/v1/orders", m...)
	transaction.GET("", h.FindAll, h.EchoRoute.Authentication)
	transaction.POST("", h.CreateOrder, h.EchoRoute.Authentication)
	transaction.POST("/accept", h.AcceptOrder, h.EchoRoute.Authentication)
}

// FindAll godoc
// @Summary Get all orders
// @Description Get all orders for the authenticated user
// @Tags Orders
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utilities.ResponseRequest
// @Failure 400 {object} utilities.ResponseRequest
// @Router /v1/orders [get]
func (h *Handler) FindAll(c echo.Context) error {
	var reqData = new(dto.FindAllRequest)

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}

	reqData.UserID = data.UserID
	reqData.RoleID = data.Role

	resp, err := h.Logic.FindAll(c.Request().Context(), reqData)
	if err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: err,
		})
	}

	return utilities.Response(c, &utilities.ResponseRequest{
		Code:   http.StatusOK,
		Status: static.Success,
		Data:   resp,
	})
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderRequest true "Create Order Request"
// @Security Bearer
// @Success 200 {object} utilities.ResponseRequest
// @Failure 400 {object} utilities.ResponseRequest
// @Router /v1/orders [post]
func (h *Handler) CreateOrder(c echo.Context) error {
	var reqData = new(dto.CreateOrderRequest)

	if err := c.Bind(reqData); err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.BadRequest), http.StatusBadRequest),
		})
	}

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}

	reqData.UserID = data.UserID
	reqData.RoleID = data.Role

	tx := h.Db.Gorm.Begin()
	if err := h.Logic.CreateOrder(c.Request().Context(), reqData, tx); err != nil {
		h.Logger.Error(err)
		defer func() {
			tx.Rollback()
		}()
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: err,
		})
	}
	tx.Commit()

	return utilities.Response(c, &utilities.ResponseRequest{
		Code:   http.StatusOK,
		Status: static.Success,
	})
}

// AcceptOrder godoc
// @Summary Accept an order
// @Description Accept an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body dto.AcceptOrderRequest true "Accept Order Request"
// @Security Bearer
// @Success 200 {object} utilities.ResponseRequest
// @Failure 400 {object} utilities.ResponseRequest
// @Router /v1/orders/accept [post]
func (h *Handler) AcceptOrder(c echo.Context) error {
	var reqData = new(dto.AcceptOrderRequest)

	if err := c.Bind(reqData); err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.BadRequest), http.StatusBadRequest),
		})
	}

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}

	reqData.DriverID = data.UserID
	reqData.RoleID = data.Role

	tx := h.Db.Gorm.Begin()
	if err := h.Logic.AcceptOrder(c.Request().Context(), reqData, tx); err != nil {
		h.Logger.Error(err)
		defer func() {
			tx.Rollback()
		}()
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: err,
		})
	}
	tx.Commit()

	return utilities.Response(c, &utilities.ResponseRequest{
		Code:   http.StatusOK,
		Status: static.Success,
	})
}
