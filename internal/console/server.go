package console

import (
	"log"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tubagusmf/ecommerce-user-product-service/db"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/config"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/repository"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/usecase"

	orderpb "github.com/tubagusmf/ecommerce-user-product-service/pb/order"
	productpb "github.com/tubagusmf/ecommerce-user-product-service/pb/product"
	userpb "github.com/tubagusmf/ecommerce-user-product-service/pb/user"

	handlerGrpc "github.com/tubagusmf/ecommerce-user-product-service/internal/delivery/grpc"
	handlerHttp "github.com/tubagusmf/ecommerce-user-product-service/internal/delivery/http"
)

var startServeCmd = &cobra.Command{
	Use:   "httpsrv",
	Short: "Start the HTTP and gRPC servers",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadWithViper()
		dbConn := db.NewPostgres()
		sqlDB, err := dbConn.DB()
		if err != nil {
			log.Fatalf("Failed to get SQL DB from Gorm: %v", err)
		}
		defer sqlDB.Close()

		// Setup repositories
		userRepo := repository.NewUserRepo(dbConn)
		productRepo := repository.NewProductRepo(dbConn)
		categoryRepo := repository.NewCategoryRepo(dbConn)
		orderRepo := repository.NewOrderRepo(dbConn)

		// Setup gRPC connections
		userConn, err := grpc.Dial("user-service:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to User Service: %v", err)
		}
		defer userConn.Close()
		userClient := userpb.NewUserServiceClient(userConn)

		productConn, err := grpc.Dial("product-service:5002", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to Product Service: %v", err)
		}
		defer productConn.Close()
		productClient := productpb.NewProductServiceClient(productConn)

		orderConn, err := grpc.Dial("order-service:5003", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to Order Service: %v", err)
		}
		defer orderConn.Close()
		orderClient := orderpb.NewOrderServiceClient(orderConn)

		// Setup usecases
		userUsecase := usecase.NewUserUsecase(userRepo, userClient)
		productUsecase := usecase.NewProductUsecase(productRepo, productClient)
		categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
		orderUsecase := usecase.NewOrderUsecase(orderRepo, productRepo, orderClient)

		quitCh := make(chan bool, 1)

		// Start HTTP server
		go func() {
			e := echo.New()
			e.GET("/ping", func(c echo.Context) error {
				return c.String(http.StatusOK, "pong!")
			})
			handlerHttp.NewUserHandler(e, userUsecase)
			handlerHttp.NewProductHandler(e, productUsecase)
			handlerHttp.NewCategoryHandler(e, categoryUsecase)
			handlerHttp.NewOrderHandler(e, orderUsecase)

			log.Println("Starting HTTP server on port 3000...")
			if err := e.Start(":3000"); err != nil {
				logrus.Fatalf("Failed to start HTTP server: %v", err)
			}
		}()

		// Start gRPC server
		go func() {
			grpcServer := grpc.NewServer()
			grpcUserHandler := handlerGrpc.NewUsergRPCHandler(userUsecase)
			grpcOrderHandler := handlerGrpc.NewOrdergRPCHandler(orderUsecase)
			grpcProductHandler := handlerGrpc.NewProductgRPCHandler(productUsecase)
			userpb.RegisterUserServiceServer(grpcServer, grpcUserHandler)
			orderpb.RegisterOrderServiceServer(grpcServer, grpcOrderHandler)
			productpb.RegisterProductServiceServer(grpcServer, grpcProductHandler)

			lis, err := net.Listen("tcp", ":5001")
			if err != nil {
				log.Fatalf("Failed to create gRPC listener: %v", err)
			}

			log.Println("gRPC server is running on port: 5001")
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalf("Failed to start gRPC server: %v", err)
			}
		}()

		<-quitCh
	},
}

func init() {
	rootCmd.AddCommand(startServeCmd)
}
