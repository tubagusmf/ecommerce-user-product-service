package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
)

type CategoryHandler struct {
	categoryUsecase model.ICategoryUsecase
}

func NewCategoryHandler(e *echo.Echo, categoryUsecase model.ICategoryUsecase) {
	handler := &CategoryHandler{
		categoryUsecase: categoryUsecase,
	}

	routeCategory := e.Group("/v1/categories")
	routeCategory.GET("", handler.FindAll, AuthMiddleware)
	routeCategory.GET("/:id", handler.FindById, AuthMiddleware)
	routeCategory.POST("/create", handler.Create, AuthMiddleware)
	routeCategory.PUT("/update/:id", handler.Update, AuthMiddleware)
	routeCategory.DELETE("/delete/:id", handler.Delete, AuthMiddleware)
}

func (h *CategoryHandler) FindAll(c echo.Context) error {
	categories, err := h.categoryUsecase.FindAll(c.Request().Context(), model.Category{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   categories,
	})
}

func (h *CategoryHandler) FindById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	category, err := h.categoryUsecase.FindById(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if category == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   category,
	})
}

func (h *CategoryHandler) Create(c echo.Context) error {
	var body model.CreateCategoryInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err := h.categoryUsecase.Create(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: "Category created successfully",
	})
}

func (h *CategoryHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var body model.UpdateCategoryInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err = h.categoryUsecase.Update(c.Request().Context(), id, body)
	if err != nil {
		if err.Error() == "category not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Category not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Category updated successfully",
	})
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = h.categoryUsecase.Delete(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "category not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Category not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Category deleted successfully",
	})
}
