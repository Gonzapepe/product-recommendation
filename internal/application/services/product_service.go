package services

import "backend-challenge/internal/domain/entities"


type ProductService interface {
    GetRecommendations(productID string) ([]*Recommendation, error)
    ComputeFeatureVectors() map[string]map[string]float64
    GetProductByID(id string) (*entities.Product, error)
    GetPaginatedProducts(offset, limit int) ([]*entities.Product, error)
    GetAllProducts() ([]*entities.Product, error)
    CreateProduct(product *entities.Product) error
    UpdateProduct(product *entities.Product) error
    DeleteProduct(id string) error
}