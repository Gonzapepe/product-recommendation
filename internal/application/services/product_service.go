package services

import "backend-challenge/internal/domain/entities"


type ProductService interface {
    GetProductByID(id string) (*entities.Product, error)
    GetPaginatedProducts(offset, limit int) ([]*entities.Product, error)
    GetAllProducts() ([]*entities.Product, error)
    CreateProduct(product *entities.Product) error
    UpdateProduct(product *entities.Product) error
    DeleteProduct(id string) error
}