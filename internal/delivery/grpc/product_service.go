package grpc

import (
	"context"
	"log"

	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
	pb "github.com/tubagusmf/ecommerce-user-product-service/pb/product_service"
)

type ProductgRPCHandler struct {
	pb.UnimplementedProductServiceServer
	productUsecase model.IProductUsecase
}

func NewProductgRPCHandler(productUsecase model.IProductUsecase) pb.ProductServiceServer {
	return &ProductgRPCHandler{productUsecase: productUsecase}
}

func (h *ProductgRPCHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := h.productUsecase.FindById(ctx, req.ProductId)
	if err != nil {
		log.Println("Error fetching product:", err)
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			ProductId:   product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		},
	}, nil
}

func (h *ProductgRPCHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := h.productUsecase.FindAll(ctx, model.FindAllParam{})
	if err != nil {
		log.Println("Error fetching products:", err)
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, &pb.Product{
			ProductId:   product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		})
	}

	return &pb.ListProductsResponse{Products: pbProducts}, nil
}

func (h *ProductgRPCHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	createdProduct, err := h.productUsecase.Create(ctx, model.CreateProductInput{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	})

	if err != nil {
		log.Println("Error creating product:", err)
		return nil, err
	}

	return &pb.CreateProductResponse{
		Product: &pb.Product{
			ProductId:   createdProduct.ID,
			Name:        createdProduct.Name,
			Description: createdProduct.Description,
			Price:       createdProduct.Price,
			Stock:       createdProduct.Stock,
		},
	}, nil
}

func (h *ProductgRPCHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	input := model.UpdateProductInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	updatedProduct, err := h.productUsecase.Update(ctx, req.ProductId, input)
	if err != nil {
		log.Println("Error updating product:", err)
		return nil, err
	}

	return &pb.UpdateProductResponse{
		Product: &pb.Product{
			ProductId:   updatedProduct.ID,
			Name:        updatedProduct.Name,
			Description: updatedProduct.Description,
			Price:       updatedProduct.Price,
			Stock:       updatedProduct.Stock,
		},
	}, nil
}

func (h *ProductgRPCHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := h.productUsecase.Delete(ctx, req.ProductId)
	if err != nil {
		log.Println("Error deleting product:", err)
		return nil, err
	}

	return &pb.DeleteProductResponse{Success: true}, nil
}
