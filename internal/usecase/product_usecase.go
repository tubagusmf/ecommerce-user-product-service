package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/helper"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
	"github.com/tubagusmf/ecommerce-user-product-service/pb/product_service"
)

type ProductUsecase struct {
	productRepo   model.IProductRepository
	productClient product_service.ProductServiceClient
}

func NewProductUsecase(
	productRepo model.IProductRepository,
	productClient product_service.ProductServiceClient,
) model.IProductUsecase {
	return &ProductUsecase{
		productRepo:   productRepo,
		productClient: productClient,
	}
}

func (u *ProductUsecase) FindAll(ctx context.Context, filter model.FindAllParam) ([]*model.Product, error) {
	log := logrus.WithFields(logrus.Fields{
		"filter": filter,
	})

	products, err := u.productRepo.FindAll(ctx, filter)
	if err != nil {
		log.Error("Failed to fetch products: ", err)
		return nil, err
	}

	return products, nil
}

func (u *ProductUsecase) FindById(ctx context.Context, id int64) (*model.Product, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	product, err := u.productRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to fetch product by ID: ", err)
		return nil, err
	}

	return product, nil
}

func (u *ProductUsecase) Create(ctx context.Context, in model.CreateProductInput) (model.Product, error) {
	log := logrus.WithFields(logrus.Fields{
		"in": in,
	})

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error:", err)
		return model.Product{}, err
	}

	if in.Name == "" || in.Price <= 0 || in.Stock < 0 || in.ImageUrl == "" {
		return model.Product{}, errors.New("invalid product data")
	}

	product := model.Product{
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Stock:       in.Stock,
		CategoryID:  in.CategoryID,
		ImageUrl:    in.ImageUrl,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		log.Error("Failed to create product: ", err)
		return model.Product{}, err
	}

	return product, nil
}

func (u *ProductUsecase) Update(ctx context.Context, id int64, in model.UpdateProductInput) (*model.Product, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
		"in": in,
	})

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error:", err)
		return &model.Product{}, err
	}

	existingProduct, err := u.productRepo.FindById(ctx, id)
	if err != nil {
		return &model.Product{}, err
	}

	existingProduct.Name = in.Name
	existingProduct.Description = in.Description
	existingProduct.Price = in.Price
	existingProduct.Stock = in.Stock
	existingProduct.CategoryID = in.CategoryID
	existingProduct.ImageUrl = in.ImageUrl
	existingProduct.UpdatedAt = time.Now()

	if err := u.productRepo.Update(ctx, *existingProduct); err != nil {
		log.Error("Failed to update product: ", err)
		return &model.Product{}, err
	}

	return existingProduct, nil
}

func (u *ProductUsecase) Delete(ctx context.Context, id int64) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	product, err := u.productRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to find product for deletion: ", err)
		return err
	}

	if product == nil {
		log.Error("Product not found")
		return errors.New("product not found")
	}

	if product.DeletedAt != nil {
		log.Error("Product already deleted")
		return errors.New("product already deleted")
	}

	err = u.productRepo.Delete(ctx, id)
	if err != nil {
		log.Error("Failed to delete product: ", err)
		return err
	}

	log.Info("Successfully deleted product with ID: ", id)
	return nil
}
