package repositories

import "backend-challenge/internal/domain/entities"

// CategoryRepository is the port for interacting with categories in the domain layer
type CategoryRepository interface {
	GetByID(id string) (*entities.Category, error)
	Create(category *entities.Category) error
	Update(category *entities.Category) error
	Delete(id string) error
	GetAll() ([]*entities.Category, error)
}
