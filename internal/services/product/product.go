package product

import (
	"local/gorm-example/internal/services/product/models"
	"log/slog"

	"gorm.io/gorm"
)

type productService struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewService(db *gorm.DB, log *slog.Logger) Service {
	return &productService{db: db, log: log}
}

func (s *productService) List() ([]models.Product, error) {
	const fn = "internal.product.List"

	s.log.Info("Finding products", slog.String("fn", fn))

	var products []models.Product

	if err := s.db.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *productService) GetOne(id uint) (*models.Product, error) {
	var product models.Product

	if result := s.db.First(&product, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (s *productService) Create(product *models.Product) error {
	if result := s.db.Create(&product); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *productService) Delete(id uint) error {
	if result := s.db.Delete(&models.Product{}, id); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *productService) Update(id uint, product *models.Product) error {
	if result := s.db.Where("id = ?", id).Updates(product); result.Error != nil {
		return result.Error
	}

	return nil
}
