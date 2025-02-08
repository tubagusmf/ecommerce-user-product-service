package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase model.IUserUsecase
}

func NewUserHandler(e *echo.Echo, userUsecase model.IUserUsecase) {
	handlers := &UserHandler{
		userUsecase: userUsecase,
	}

	routeUser := e.Group("v1/auth")
	routeUser.POST("/login", handlers.Login)
	routeUser.GET("/user/:id", handlers.FindById, AuthMiddleware)
	routeUser.GET("/users", handlers.FindAll, AuthMiddleware)
	routeUser.POST("/register", handlers.Create)
	routeUser.PUT("/user/update/:id", handlers.Update, AuthMiddleware)
	routeUser.DELETE("/user/delete/:id", handlers.Delete, AuthMiddleware)
}

func (handler *UserHandler) Login(c echo.Context) error {
	var body model.LoginInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accessToken, err := handler.userUsecase.Login(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:      http.StatusOK,
		Message:     "Success Login",
		AccessToken: accessToken,
	})
}

func (handler *UserHandler) FindById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	claim, ok := c.Request().Context().Value(model.BearerAuthKey).(model.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	log.Printf("Authenticated User ID: %d", claim.UserID)

	user, err := handler.userUsecase.FindById(c.Request().Context(), int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   user,
	})
}

func (handler *UserHandler) FindAll(c echo.Context) error {
	users, err := handler.userUsecase.FindAll(c.Request().Context(), model.User{})
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   users,
	})
}

func (handler *UserHandler) Create(c echo.Context) error {
	var body model.CreateUserInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accessToken, err := handler.userUsecase.Create(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:      http.StatusOK,
		Message:     "Success Register",
		AccessToken: accessToken,
	})
}

func (handler *UserHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	claim, ok := c.Request().Context().Value(model.BearerAuthKey).(model.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	log.Printf("Authenticated User ID: %d", claim.UserID)

	var body model.UpdateUserInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = handler.userUsecase.Update(c.Request().Context(), id, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "User updated successfully",
		Data:    body,
	})
}

func (handler *UserHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	claim, ok := c.Request().Context().Value(model.BearerAuthKey).(model.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	log.Printf("Authenticated User ID: %d", claim.UserID)

	err = handler.userUsecase.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "User deleted successfully",
	})
}
