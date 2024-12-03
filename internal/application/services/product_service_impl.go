package services

import (
	"backend-challenge/internal/domain/entities"
	"backend-challenge/internal/domain/repositories"
)

type productService struct {
	repo repositories.ProductRepository
    recommender *RecommendationService
}

func NewProductService(repo repositories.ProductRepository, recommender *RecommendationService) ProductService {
	return &productService{repo: repo, recommender: recommender}
}

func (s *productService) ComputeFeatureVectors() map[string]map[string]float64 {
    products, _ := s.repo.GetAll()

    featureVectors := make(map[string]map[string]float64)

    for _, product := range products {
        featureVectors[product.ID.Hex()] = s.recommender.ExtractFeatureVector(*product)
    }

    return featureVectors
}

func (s *productService) GetProductByID(id string) (*entities.Product, error) {
    return s.repo.GetByID(id)
}

func (s *productService) GetPaginatedProducts(offset, limit int) ([]*entities.Product, error) {
    return s.repo.GetPaginated(offset, limit)
}

func (s *productService) GetAllProducts() ([]*entities.Product, error) {
    return s.repo.GetAll()
}

func (s *productService) CreateProduct(product *entities.Product) error {
    return s.repo.Create(product)
}

func (s *productService) UpdateProduct(product *entities.Product) error {
    return s.repo.Update(product)
}

func (s *productService) DeleteProduct(id string) error {
    return s.repo.Delete(id)
}