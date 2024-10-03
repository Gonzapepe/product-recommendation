package services

import "backend-challenge/internal/domain/entities"

type CategoryService interface {
	GetCategoryByID(id string) (*entities.Category, error)
	GetAllCategories() ([]*entities.Category, error)
	CreateCategory(category *entities.Category) error
	UpdateCategory(category *entities.Category) error
	DeleteCategory(id string) error
}