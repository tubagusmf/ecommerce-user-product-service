package http

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/tubagusmf/ecommerce-user-product-service/internal/helper"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		splitAuth := strings.Split(authHeader, " ")
		if len(splitAuth) != 2 || splitAuth[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		accessToken := splitAuth[1]
		log.Println(accessToken)

		var claim model.CustomClaims
		err := helper.DecodeToken(accessToken, &claim)
		log.Println(claim)
		if err != nil {
			log.Println(claim)
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}
		ctx := context.WithValue(c.Request().Context(), model.BearerAuthKey, claim)
		req := c.Request().WithContext(ctx)
		c.SetRequest(req)

		return next(c)
	}
}
