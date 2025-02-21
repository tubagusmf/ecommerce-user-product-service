package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/helper"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
)

type OrderUsecase struct {
	orderRepo   model.IOrderRepository
	productRepo model.IProductRepository
}

func NewOrderUsecase(orderRepo model.IOrderRepository, productRepo model.IProductRepository) model.IOrderUsecase {
	return &OrderUsecase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (u *OrderUsecase) FindAll(ctx context.Context, userID int64) ([]*model.Order, error) {
	log := logrus.WithFields(logrus.Fields{
		"userID": userID,
	})

	if userID == 0 {
		log.Error("Invalid user ID")
		return nil, errors.New("invalid user ID")
	}

	orders, err := u.orderRepo.FindAll(ctx, userID)
	if err != nil {
		log.Error("Failed to fetch orders: ", err)
		return nil, err
	}

	return orders, nil
}

func (u *OrderUsecase) FindById(ctx context.Context, id string) (*model.Order, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	if id == "" {
		log.Error("Invalid order ID")
		return nil, errors.New("invalid order ID")
	}

	order, err := u.orderRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to fetch order: ", err)
		return nil, err
	}

	return order, nil
}

func (u *OrderUsecase) ListByUserID(ctx context.Context, userID int64) ([]*model.Order, error) {
	log := logrus.WithFields(logrus.Fields{
		"user_id": userID,
	})

	log.Info("Fetching orders for user...")

	orders, err := u.orderRepo.FindAll(ctx, userID)
	if err != nil {
		log.WithError(err).Error("Failed to fetch orders")
		return nil, err
	}

	log.Infof("Found %d orders for user", len(orders))
	return orders, nil
}

func (u *OrderUsecase) Create(ctx context.Context, in model.CreateOrderInput) (*model.Order, error) {
	log := logrus.WithFields(logrus.Fields{
		"in": in,
	})

	if len(in.OrderItems) == 0 {
		return nil, errors.New("order items cannot be empty")
	}

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error:", err)
		return nil, err
	}

	order := model.Order{
		UserID:      in.UserID,
		TotalAmount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, item := range in.OrderItems {
		price, err := u.getProductPrice(ctx, item.ProductID)
		if err != nil {
			log.Error("Failed to get product price: ", err)
			return nil, err
		}

		order.TotalAmount += price * float64(item.Quantity)

		order.OrderItems = append(order.OrderItems, model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	if err := u.orderRepo.SaveOrder(ctx, &order); err != nil {
		log.Error("Failed to save order: ", err)
		return nil, err
	}

	return &order, nil
}

func (u *OrderUsecase) Update(ctx context.Context, order *model.Order) error {
	if order == nil || order.ID == "" {
		logrus.WithContext(ctx).Error("Invalid order: nil or empty ID")
		return errors.New("invalid order: order is nil or has empty ID")
	}

	log := logrus.WithFields(logrus.Fields{
		"order_id": order.ID,
		"user_id":  order.UserID,
	})

	log.Info("Updating order...")

	err := u.orderRepo.SaveOrder(ctx, order)
	if err != nil {
		log.WithError(err).Error("Failed to update order")
		return err
	}

	log.Info("Order successfully updated")
	return nil
}

func (u *OrderUsecase) Delete(ctx context.Context, id string) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	order, err := u.orderRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to find order for deletion: ", err)
		return err
	}

	if order == nil {
		log.Error("Order not found")
		return errors.New("order not found")
	}

	if order.DeletedAt != nil {
		log.Error("Order already deleted")
		return errors.New("order already deleted")
	}

	err = u.orderRepo.Delete(ctx, id)
	if err != nil {
		log.Error("Failed to delete order: ", err)
		return err
	}

	log.Info("Successfully deleted order with ID: ", id)
	return nil
}

func (u *OrderUsecase) getProductPrice(ctx context.Context, productID int64) (float64, error) {
	var price float64
	err := u.productRepo.GetPriceByID(ctx, productID, &price)
	if err != nil {
		return 0, err
	}
	return price, nil
}
