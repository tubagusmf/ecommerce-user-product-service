package console

import (
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tubagusmf/ecommerce-user-product-service/db"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/config"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/repository"
	"github.com/tubagusmf/ecommerce-user-product-service/internal/usecase"

	orderpb "github.com/tubagusmf/ecommerce-user-product-service/pb/order_service"
	productpb "github.com/tubagusmf/ecommerce-user-product-service/pb/product_service"
	userpb "github.com/tubagusmf/ecommerce-user-product-service/pb/user_service"

	handlerHttp "github.com/tubagusmf/ecommerce-user-product-service/internal/delivery/http"
)

func init() {
	rootCmd.AddCommand(serverCMD)
}

var serverCMD = &cobra.Command{
	Use:   "httpsrv",
	Short: "Start HTTP and gRPC server",
	Long:  "Start the HTTP server to handle requests and establish gRPC connections to services.",
	Run:   startServer,
}

func startServer(cmd *cobra.Command, args []string) {
	config.LoadWithViper()

	// Initialize PostgreSQL connection
	postgresDB := db.NewPostgres()
	sqlDB, err := postgresDB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB from Gorm: %v", err)
	}
	defer sqlDB.Close()

	// Setup repositories
	userRepo := repository.NewUserRepo(postgresDB)
	productRepo := repository.NewProductRepo(postgresDB)
	categoryRepo := repository.NewCategoryRepo(postgresDB)
	orderRepo := repository.NewOrderRepo(postgresDB)

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

	// Initialize Echo HTTP server
	e := echo.New()

	// Register HTTP handlers
	handlerHttp.NewUserHandler(e, userUsecase)
	handlerHttp.NewProductHandler(e, productUsecase)
	handlerHttp.NewCategoryHandler(e, categoryUsecase)
	handlerHttp.NewOrderHandler(e, orderUsecase)

	// Start HTTP server and gRPC server concurrently
	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Println("Starting HTTP server on port 3000...")
		errCh <- e.Start(":3000")
	}()

	go func() {
		defer wg.Done()
		<-errCh
	}()

	wg.Wait()

	if err := <-errCh; err != nil {
		if err != http.ErrServerClosed {
			logrus.Errorf("HTTP server error: %v", err)
		}
	}
}
