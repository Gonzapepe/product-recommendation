package repositories

import "backend-challenge/internal/domain/entities"

type ProductRepository interface {
	GetByID(id string) (*entities.Product, error)
	GetPaginated(offset, limit int) ([]*entities.Product, error)
	GetAll() ([]*entities.Product, error)
	Create(product *entities.Product) error
	Update(product *entities.Product) error
	Delete(id string) error
}