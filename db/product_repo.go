package db

import (
	"errors"
	"log"
	"github.com/techagentng/ecommerce-api/models"
	"gorm.io/gorm"
)

// ProductRepository interface defines the methods for product-related database operations
type ProductRepository interface {
	CreateProduct(product *models.Product) (*models.Product, error)
	FindProductByID(id uint) (*models.Product, error)
	FindAllProducts() ([]*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint) error
}

// productRepo struct holds the database connection
type productRepo struct {
	DB *gorm.DB
}

// NewProductRepo creates a new instance of ProductRepository
func NewProductRepo(db *GormDB) ProductRepository {
	return &productRepo{db.DB}
}

// CreateProduct inserts a new product into the database
func (p *productRepo) CreateProduct(product *models.Product) (*models.Product, error) {
	if err := p.DB.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

// FindProductByID retrieves a product by its ID
func (p *productRepo) FindProductByID(id uint) (*models.Product, error) {
    var product models.Product
    log.Printf("Attempting to find product with ID: %d", id)

    if err := p.DB.First(&product, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            log.Printf("Product with ID: %d not found", id)
            return nil, nil
        }
        log.Printf("Error retrieving product with ID: %d, Error: %v", id, err)
        return nil, err
    }

    log.Printf("Product found: %+v", product)
    return &product, nil
}

// FindAllProducts retrieves all products from the database
func (p *productRepo) FindAllProducts() ([]*models.Product, error) {
	var products []*models.Product
	if err := p.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateProduct updates an existing product in the database
func (p *productRepo) UpdateProduct(product *models.Product) error {
	return p.DB.Save(product).Error
}

// DeleteProduct removes a product from the database
func (p *productRepo) DeleteProduct(id uint) error {
	return p.DB.Delete(&models.Product{}, id).Error
}
