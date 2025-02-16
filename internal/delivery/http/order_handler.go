package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
)

type OrderHandler struct {
	orderUsecase model.IOrderUsecase
}

func NewOrderHandler(e *echo.Echo, orderUsecase model.IOrderUsecase) {
	handler := &OrderHandler{
		orderUsecase: orderUsecase,
	}

	routeOrder := e.Group("v1/orders")
	routeOrder.GET("", handler.FindAll, AuthMiddleware)
	routeOrder.GET("/:id", handler.FindById, AuthMiddleware)
	routeOrder.POST("/create", handler.Create, AuthMiddleware)
	routeOrder.DELETE("/delete/:id", handler.Delete, AuthMiddleware)
}

func (handler *OrderHandler) FindAll(c echo.Context) error {
	userIDStr := c.QueryParam("user_id")

	if userIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is required")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	orders, err := handler.orderUsecase.FindAll(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   orders,
	})
}

func (handler *OrderHandler) FindById(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Order ID is required")
	}

	order, err := handler.orderUsecase.FindById(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "order not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Order not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   order,
	})
}

func (handler *OrderHandler) Create(c echo.Context) error {
	var body model.CreateOrderInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := handler.orderUsecase.Create(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Order created successfully",
	})
}

func (handler *OrderHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	err := handler.orderUsecase.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Order deleted successfully",
	})
}
