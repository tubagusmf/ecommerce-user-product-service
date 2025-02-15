package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/helper"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
)

type CategoryUsecase struct {
	categoryRepo model.ICategoryRepository
}

func NewCategoryUsecase(categoryRepo model.ICategoryRepository) model.ICategoryUsecase {
	return &CategoryUsecase{categoryRepo: categoryRepo}
}

func (u *CategoryUsecase) FindAll(ctx context.Context, category model.Category) ([]*model.Category, error) {
	log := logrus.WithFields(logrus.Fields{
		"category": category,
	})

	categories, err := u.categoryRepo.FindAll(ctx, category)
	if err != nil {
		log.Error("Failed to fetch categories: ", err)
		return nil, err
	}

	return categories, nil
}

func (u *CategoryUsecase) FindById(ctx context.Context, id int64) (*model.Category, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	category, err := u.categoryRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to fetch category by ID: ", err)
		return nil, err
	}

	return category, nil
}

func (u *CategoryUsecase) Create(ctx context.Context, in model.CreateCategoryInput) error {
	log := logrus.WithFields(logrus.Fields{
		"in": in,
	})

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error:", err)
		return err
	}

	category := model.Category{
		Name: in.Name,
	}

	if err := u.categoryRepo.Create(ctx, category); err != nil {
		log.Error("Failed to create category: ", err)
		return err
	}

	return nil
}

func (u *CategoryUsecase) Update(ctx context.Context, id int64, in model.UpdateCategoryInput) error {
	log := logrus.WithFields(logrus.Fields{
		"id":   id,
		"in":   in,
		"name": in.Name,
	})

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error:", err)
		return err
	}

	existingCategory, err := u.categoryRepo.FindById(ctx, id)
	if err != nil {
		return err
	}

	existingCategory.Name = in.Name
	existingCategory.UpdatedAt = time.Now()

	if err := u.categoryRepo.Update(ctx, *existingCategory); err != nil {
		log.Error("Failed to update category: ", err)
		return err
	}

	return nil
}

func (u *CategoryUsecase) Delete(ctx context.Context, id int64) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	category, err := u.categoryRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to find category for deletion: ", err)
		return err
	}

	if category == nil {
		log.Error("Category not found")
		return errors.New("category not found")
	}

	if category.DeletedAt != nil {
		log.Error("Category already deleted")
		return errors.New("category already deleted")
	}

	if err := u.categoryRepo.Delete(ctx, id); err != nil {
		log.Error("Failed to delete category: ", err)
		return err
	}

	log.Info("Successfully deleted category with ID: ", id)

	return nil
}
