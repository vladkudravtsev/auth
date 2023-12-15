package api

import "local/gorm-example/internal/services/product/models"

type ProductResponse struct {
	ID    uint   `json:"id"`
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

type ProductParams struct {
	ID uint `uri:"id" binding:"required"`
}

func ToProductResponse(product models.Product) *ProductResponse {
	return &ProductResponse{
		ID:    product.ID,
		Price: product.Price,
		Code:  product.Code,
	}
}
