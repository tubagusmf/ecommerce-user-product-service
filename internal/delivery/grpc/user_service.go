package grpc

import (
	"context"
	"log"

	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
	pb "github.com/tubagusmf/ecommerce-user-product-service/pb/user"
)

// UsergRPCHandler mengimplementasikan UserServiceServer
type UsergRPCHandler struct {
	pb.UnimplementedUserServiceServer
	userUsecase model.IUserUsecase
}

// NewUsergRPCHandler adalah constructor untuk UsergRPCHandler
func NewUsergRPCHandler(userUsecase model.IUserUsecase) pb.UserServiceServer {
	return &UsergRPCHandler{userUsecase: userUsecase}
}

// GetUser menghandle request GetUser dan mengembalikan data user sesuai dengan protokol gRPC
func (h *UsergRPCHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Ambil data user dari usecase
	user, err := h.userUsecase.FindById(ctx, req.GetUserId())
	if err != nil {
		log.Println("Error fetching user:", err)
		return nil, err
	}

	// Konversi user ke protobuf response
	response := &pb.GetUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}

	return response, nil
}
