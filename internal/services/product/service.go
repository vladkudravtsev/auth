package product

import "local/gorm-example/internal/services/product/models"

type Service interface {
	List() ([]models.Product, error)
	GetOne(id uint) (*models.Product, error)
	Create(product *models.Product) error
	Delete(id uint) error
	Update(id uint, product *models.Product) error
}
