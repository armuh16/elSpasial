package router

import (
	"context"
	"errors"
	"github.com/elspasial/config"
	"github.com/elspasial/model"
	"github.com/elspasial/package/jwt"
	"github.com/elspasial/static"
	"github.com/elspasial/utilities"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

//func (r *Router) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		AuthorizationHeader := c.Request().Header.Get("Authorization")
//		r.Logger.Error("Authorization header missing")
//		Authorization := strings.Split(AuthorizationHeader, " ")
//		if len(Authorization) > 1 {
//			result, err := jwt.ParseClaim(Authorization[1], config.Get().Auth.Secret)
//			if err != nil {
//				r.Logger.Error(err.Error())
//				return utilities.Response(c, &utilities.ResponseRequest{
//					Code:  http.StatusUnauthorized,
//					Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
//				})
//			}
//
//			ctx := c.Request().Context()
//
//			// Check exist user
//			userDetail, err := r.userRepo.Find(ctx, &model.Users{
//				ID: result.Data.UserID,
//			})
//			if err != nil {
//				r.Logger.Error(err.Error())
//				return utilities.Response(c, &utilities.ResponseRequest{
//					Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
//				})
//			}
//
//			ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
//				UserID: userDetail.ID,
//				Role:   userDetail.RoleID,
//			})
//
//			c.SetRequest(c.Request().WithContext(ctx))
//
//		} else {
//			return utilities.Response(c, &utilities.ResponseRequest{
//				Code:  http.StatusUnauthorized,
//				Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
//			})
//		}
//		return next(c)
//	}
//}

// Authentication middleware
//func (r *Router) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		AuthorizationHeader := c.Request().Header.Get("Authorization")
//		r.Logger.Info("Authorization Header: ", AuthorizationHeader)
//
//		if AuthorizationHeader == "" {
//			r.Logger.Error("Authorization header missing")
//			return utilities.Response(c, &utilities.ResponseRequest{
//				Code:  http.StatusUnauthorized,
//				Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
//			})
//		}
//
//		Authorization := strings.Split(AuthorizationHeader, " ")
//		if len(Authorization) != 2 || Authorization[0] != "Bearer" {
//			r.Logger.Error("Authorization header format invalid")
//			return utilities.Response(c, &utilities.ResponseRequest{
//				Code:  http.StatusUnauthorized,
//				Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
//			})
//		}
//
//		token := Authorization[1]
//		r.Logger.Info("Token received: ", token)
//
//		result, err := jwt.ParseClaim(token, config.Get().Auth.Secret)
//		if err != nil {
//			r.Logger.Error("Token parsing failed: ", err.Error())
//			return utilities.Response(c, &utilities.ResponseRequest{
//				Code:  http.StatusUnauthorized,
//				Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
//			})
//		}
//
//		r.Logger.Info("Token valid, user ID: ", result.Data.UserID)
//
//		ctx := c.Request().Context()
//
//		// Check exist user
//		userDetail, err := r.userRepo.Find(ctx, &model.Users{
//			ID: result.Data.UserID,
//		})
//		if err != nil {
//			r.Logger.Error("User not found: ", err.Error())
//			return utilities.Response(c, &utilities.ResponseRequest{
//				Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
//			})
//		}
//
//		r.Logger.Info("User found: ", userDetail.Email)
//
//		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
//			UserID: userDetail.ID,
//			Role:   userDetail.RoleID,
//		})
//
//		c.SetRequest(c.Request().WithContext(ctx))
//
//		return next(c)
//	}
//}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// Authentication middleware
func (r *Router) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		AuthorizationHeader := c.Request().Header.Get("Authorization")
		if AuthorizationHeader == "" {
			r.Logger.Error("Authorization header missing")
			return utilities.Response(c, &utilities.ResponseRequest{
				Code:  http.StatusUnauthorized,
				Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
			})
		}

		Authorization := strings.Split(AuthorizationHeader, " ")
		if len(Authorization) != 2 || Authorization[0] != "Bearer" {
			r.Logger.Error("Authorization header format invalid")
			return utilities.Response(c, &utilities.ResponseRequest{
				Code:  http.StatusUnauthorized,
				Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
			})
		}

		tokenString := Authorization[1]
		r.Logger.Info("Received token: ", tokenString)

		result, err := jwt.ParseClaim(tokenString, config.Get().Auth.Secret)
		if err != nil {
			r.Logger.Error("Failed to parse JWT: ", err.Error())
			return utilities.Response(c, &utilities.ResponseRequest{
				Code:  http.StatusUnauthorized,
				Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
			})
		}

		ctx := c.Request().Context()

		// Check exist user
		userDetail, err := r.userRepo.Find(ctx, &model.Users{
			ID: result.Data.UserID,
		})
		if err != nil {
			r.Logger.Error("User not found: ", err.Error())
			return utilities.Response(c, &utilities.ResponseRequest{
				Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
			})
		}

		r.Logger.Info("Authenticated user: ", userDetail.ID)

		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: userDetail.ID,
			Role:   userDetail.RoleID,
		})

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
