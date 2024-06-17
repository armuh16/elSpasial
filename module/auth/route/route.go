package route

import (
	"errors"
	"net/http"

	"github.com/elspasial/database/postgres"
	"github.com/elspasial/module/auth/dto"
	"github.com/elspasial/module/auth/logic"
	"github.com/elspasial/package/logger"
	"github.com/elspasial/router"
	"github.com/elspasial/static"
	"github.com/elspasial/utilities"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type handler struct {
	fx.In
	Logic     logic.IAuthLogic
	EchoRoute *router.Router
	Logger    *logger.LogRus
	Db        *postgres.DB
}

func NewRoute(h handler, m ...echo.MiddlewareFunc) handler {
	h.Route(m...)
	return h
}

func (h *handler) Route(m ...echo.MiddlewareFunc) {
	auth := h.EchoRoute.Group("/v1/auth", m...)
	auth.POST("/login", h.Login)
	auth.POST("/register", h.Register)
}

// Login godoc
// @Summary User login
// @Description User login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} utilities.ResponseRequest
// @Failure 400 {object} utilities.ResponseRequest
// @Router /v1/auth/login [post]
func (h *handler) Login(c echo.Context) error {
	var reqData = new(dto.LoginRequest)

	if err := c.Bind(reqData); err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.BadRequest), http.StatusBadRequest),
		})
	}

	tx := h.Db.Gorm.Begin()
	resp, err := h.Logic.Login(c.Request().Context(), reqData, tx)
	if err != nil {
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
		Data:   resp,
	})
}

// Register godoc
// @Summary User registration
// @Description User registration
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 200 {object} utilities.ResponseRequest
// @Failure 400 {object} utilities.ResponseRequest
// @Router /v1/auth/register [post]
func (h *handler) Register(c echo.Context) error {
	var reqData = new(dto.RegisterRequest)

	if err := c.Bind(reqData); err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.BadRequest), http.StatusBadRequest),
		})
	}

	tx := h.Db.Gorm.Begin()
	if err := h.Logic.Register(c.Request().Context(), reqData, tx); err != nil {
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
