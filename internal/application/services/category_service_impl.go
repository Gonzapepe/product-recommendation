package services

import (
	"backend-challenge/internal/domain/entities"
	"backend-challenge/internal/domain/repositories"
)

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService  {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAllCategories() ([]*entities.Category, error) {
	return s.repo.GetAll()
}

func (s *categoryService) GetCategoryByID(id string) (*entities.Category, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) CreateCategory(category *entities.Category) error {
	return s.repo.Create(category)
}

func (s *categoryService) UpdateCategory(category *entities.Category) error {
	return s.repo.Update(category)
}

func (s *categoryService) DeleteCategory(id string) error {
	return s.repo.Delete(id)
}
