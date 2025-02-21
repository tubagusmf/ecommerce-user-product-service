package grpc

import (
	"context"
	"log"

	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
	pb "github.com/tubagusmf/ecommerce-user-product-service/pb/user_service"
)

type UsergRPCHandler struct {
	pb.UnimplementedUserServiceServer
	userUsecase model.IUserUsecase
}

func NewUsergRPCHandler(userUsecase model.IUserUsecase) pb.UserServiceServer {
	return &UsergRPCHandler{userUsecase: userUsecase}
}

func (h *UsergRPCHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.userUsecase.FindById(ctx, req.UserId)
	if err != nil {
		log.Println("Error fetching user:", err)
		return nil, err
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			UserId: user.ID,
			Name:   user.Name,
			Email:  user.Email,
		},
	}, nil
}
