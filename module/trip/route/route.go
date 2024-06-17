package route

import (
	"errors"
	"github.com/elspasial/enum"
	"net/http"

	"github.com/elspasial/database/postgres"
	"github.com/elspasial/module/trip/dto"
	"github.com/elspasial/module/trip/logic"
	"github.com/elspasial/package/jwt"
	"github.com/elspasial/package/logger"
	"github.com/elspasial/router"
	"github.com/elspasial/static"
	"github.com/elspasial/utilities"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type handler struct {
	fx.In
	Logic     logic.IProductLogic
	EchoRoute *router.Router
	Logger    *logger.LogRus
	Db        *postgres.DB
}

func NewRoute(h handler, m ...echo.MiddlewareFunc) handler {
	h.Route(m...)
	return h
}

func (h *handler) Route(m ...echo.MiddlewareFunc) {
	product := h.EchoRoute.Group("/v1/trip", m...)
	product.POST("", h.Create, h.EchoRoute.Authentication)
	product.GET("", h.FindAll, h.EchoRoute.Authentication)
}

// Get all trips for the authenticated user
// @Summary Get all trips
// @Description Get all trips for the authenticated user
// @Tags Trips
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} utilities.ResponseRequest
// @Failure 401 {object} utilities.ResponseRequest
// @Router /v1/trip [get]
func (h *handler) FindAll(c echo.Context) error {
	var reqData = new(dto.FindAllRequest)

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok || data.Role != enum.RoleTypeUser {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}
	reqData.UserID = data.UserID

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

// CreateRequest godoc
// @Summary Create a new trip
// @Description Create a new trip
// @Tags Trips
// @Accept json
// @Produce json
// @Param request body dto.CreateRequest true "Create Trip Request"
// @Success 200 {object} utilities.ResponseRequest
// @Failure 400 {object} utilities.ResponseRequest
// @Router /v1/trip [post]
// @Security Bearer
func (h *handler) Create(c echo.Context) error {
	var reqData = new(dto.CreateRequest)

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}

	reqData.UserID = data.UserID
	reqData.RoleID = data.Role

	if err := c.Bind(reqData); err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.BadRequest), http.StatusBadRequest),
		})
	}

	tx := h.Db.Gorm.Begin()
	if err := h.Logic.Create(c.Request().Context(), reqData, tx); err != nil {
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
