package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) model.IOrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) FindAll(ctx context.Context, userID int64) ([]*model.Order, error) {
	var orders []*model.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindById(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&order).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("order not found")
	}

	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order *model.Order) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	currentDate := time.Now().Format("20060102")

	var count int64
	if err := tx.Model(&model.Order{}).
		Where("DATE(created_at) = ?", currentDate).
		Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}
	order.ID = fmt.Sprintf("ORD-%s-%03d", currentDate, count+1)

	existingOrder := model.Order{}
	if err := tx.Where("id = ?", order.ID).First(&existingOrder).Error; err == nil {
		tx.Rollback()
		return errors.New("duplicate order detected in database")
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.OrderItems {
		var existingItem model.OrderItem
		if err := tx.Where("order_id = ? AND product_id = ?", order.ID, item.ProductID).First(&existingItem).Error; err == nil {
			existingItem.Quantity += item.Quantity
			existingItem.UpdatedAt = time.Now()
			if err := tx.Save(&existingItem).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			item.OrderID = order.ID
			item.CreatedAt = time.Now()
			item.UpdatedAt = time.Now()

			if err := tx.Create(&item).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	now := time.Now()
	err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", id).
		Update("deleted_at", now).Error

	if err != nil {
		return err
	}
	return nil
}
