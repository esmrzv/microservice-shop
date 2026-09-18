package grpc

import (
	"context"
	"errors"

	"github.com/esmrzv/product-service/internal/service"
	"github.com/esmrzv/product-service/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductServer struct {
	proto.UnimplementedProductServiceServer
	service service.ProductService
}

func NewProductServer(service service.ProductService) *ProductServer {
	return &ProductServer{service: service}
}

func (s *ProductServer) GetProduct(
	ctx context.Context,
	req *proto.GetProductRequest,
) (*proto.GetProductResponse, error) {
	productID, err := uuid.Parse(req.GetProductId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id")
	}

	product, err := s.service.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "error server error")
	}

	return &proto.GetProductResponse{
		Id:    product.ID.String(),
		Name:  product.Name,
		Price: product.Price,
	}, nil
}
