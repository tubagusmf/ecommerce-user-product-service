package grpc

import (
	"context"
	"log"

	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
	pb "github.com/tubagusmf/ecommerce-user-product-service/pb/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrdergRPCHandler struct {
	pb.UnimplementedOrderServiceServer
	orderUsecase model.IOrderUsecase
}

func NewOrdergRPCHandler(orderUsecase model.IOrderUsecase) pb.OrderServiceServer {
	return &OrdergRPCHandler{orderUsecase: orderUsecase}
}

func (h *OrdergRPCHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	orderInput := model.CreateOrderInput{
		UserID:     req.UserId,
		OrderItems: convertOrderItems(req.Items),
	}

	createdOrder, err := h.orderUsecase.Create(ctx, orderInput)
	if err != nil {
		log.Println("Error creating order:", err)
		return nil, err
	}

	return &pb.CreateOrderResponse{Order: convertOrderToPB(createdOrder)}, nil
}

func (h *OrdergRPCHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := h.orderUsecase.FindById(ctx, req.OrderId)
	if err != nil {
		log.Println("Error fetching order:", err)
		return nil, err
	}

	return &pb.GetOrderResponse{Order: convertOrderToPB(order)}, nil
}

func (h *OrdergRPCHandler) MarkOrderPaid(ctx context.Context, req *pb.MarkOrderPaidRequest) (*pb.MarkOrderPaidResponse, error) {
	order, err := h.orderUsecase.FindById(ctx, req.OrderId)
	if err != nil {
		log.Println("[ERROR] Order not found:", err)
		return &pb.MarkOrderPaidResponse{Success: false}, status.Errorf(codes.NotFound, "Order not found")
	}

	order.Status = "success"

	err = h.orderUsecase.Update(ctx, order)
	if err != nil {
		log.Println("[ERROR] Failed to update order status:", err)
		return &pb.MarkOrderPaidResponse{Success: false}, status.Errorf(codes.Internal, "Failed to update order status")
	}

	log.Printf("[INFO] Order %s marked as PAID", req.OrderId)
	return &pb.MarkOrderPaidResponse{Success: true}, nil
}

func (h *OrdergRPCHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := h.orderUsecase.ListByUserID(ctx, req.UserId)
	if err != nil {
		log.Println("Error fetching orders:", err)
		return nil, err
	}

	pbOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		pbOrders[i] = convertOrderToPB(order)
	}

	return &pb.ListOrdersResponse{Orders: pbOrders}, nil
}

func convertOrderItems(items []*pb.OrderItem) []model.CreateOrderItem {
	var orderItems []model.CreateOrderItem
	for _, item := range items {
		orderItems = append(orderItems, model.CreateOrderItem{
			ProductID: item.ProductId,
			Quantity:  int64(item.Quantity),
		})
	}
	return orderItems
}
func convertOrderToPB(order *model.Order) *pb.Order {
	pbItems := make([]*pb.OrderItem, len(order.OrderItems))
	for i, item := range order.OrderItems {
		pbItems[i] = &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return &pb.Order{
		OrderId:     order.ID,
		UserId:      order.UserID,
		Items:       pbItems,
		TotalAmount: order.TotalAmount,
	}
}
