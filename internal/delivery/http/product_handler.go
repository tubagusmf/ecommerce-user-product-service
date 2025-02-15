package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
)

type ProductHandler struct {
	productUsecase model.IProductUsecase
}

func NewProductHandler(e *echo.Echo, productUsecase model.IProductUsecase) {
	handler := &ProductHandler{
		productUsecase: productUsecase,
	}

	routeProduct := e.Group("v1/products")
	routeProduct.GET("", handler.FindAll, AuthMiddleware)
	routeProduct.GET("/:id", handler.FindById, AuthMiddleware)
	routeProduct.POST("/create", handler.Create, AuthMiddleware)
	routeProduct.PUT("/update/:id", handler.Update, AuthMiddleware)
	routeProduct.DELETE("/delete/:id", handler.Delete, AuthMiddleware)
}

func (handler *ProductHandler) FindAll(c echo.Context) error {
	var filter model.FindAllParam

	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid query parameters")
	}

	products, err := handler.productUsecase.FindAll(c.Request().Context(), filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   products,
	})
}

func (handler *ProductHandler) FindById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	product, err := handler.productUsecase.FindById(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "product not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   product,
	})
}

func (handler *ProductHandler) Create(c echo.Context) error {
	var body model.CreateProductInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := handler.productUsecase.Create(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Product created successfully",
	})
}

func (handler *ProductHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var body model.UpdateProductInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = handler.productUsecase.Update(c.Request().Context(), id, body)
	if err != nil {
		if err.Error() == "product not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Product updated successfully",
		Data:    body,
	})
}

func (handler *ProductHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = handler.productUsecase.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Product deleted successfully",
	})
}
