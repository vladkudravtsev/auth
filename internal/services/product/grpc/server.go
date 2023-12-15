package productgrpc

import (
	"context"
	productv1 "local/gorm-example/api/gen/go/product"
	"local/gorm-example/internal/services/product"
	"local/gorm-example/internal/services/product/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	productv1.UnimplementedProductServiceServer
	product product.Service
}

func RegisterServer(gRPCServer *grpc.Server, product product.Service) {
	productv1.RegisterProductServiceServer(gRPCServer, &serverAPI{product: product})
}

func (s *serverAPI) List(ctx context.Context, req *productv1.EmptyRequest) (*productv1.ListResponse, error) {
	products, err := s.product.List()

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	resp := make([]*productv1.Product, len(products))

	for _, product := range products {
		resp = append(resp, &productv1.Product{Code: product.Code, Price: uint32(product.Price)})
	}

	return &productv1.ListResponse{Products: resp}, nil
}

func (s *serverAPI) GetOne(ctx context.Context, req *productv1.GetOneRequest) (*productv1.Product, error) {
	product, err := s.product.GetOne(uint(req.GetId()))

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &productv1.Product{Code: product.Code, Price: uint32(product.Price)}, nil
}

func (s *serverAPI) Create(ctx context.Context, req *productv1.Product) (*productv1.EmptyResponse, error) {
	product := &models.Product{Code: req.GetCode(), Price: uint(req.GetPrice())}

	if err := s.product.Create(product); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &productv1.EmptyResponse{}, nil
}

func (s *serverAPI) Delete(ctx context.Context, req *productv1.DeleteRequest) (*productv1.EmptyResponse, error) {
	if err := s.product.Delete(uint(req.GetId())); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &productv1.EmptyResponse{}, nil
}

func (s *serverAPI) Update(ctx context.Context, req *productv1.UpdateRequest) (*productv1.EmptyResponse, error) {
	reqProduct := req.GetProduct()
	product := &models.Product{Code: reqProduct.GetCode(), Price: uint(reqProduct.GetPrice())}

	if err := s.product.Update(uint(req.GetId()), product); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &productv1.EmptyResponse{}, nil
}
