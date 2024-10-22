package unit

import (
	"backend-challenge/internal/adapters/persistence/repository"
	"backend-challenge/internal/domain/entities"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var testDB *mongo.Client

func TestMainCategory(m *testing.M) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	testDB = client

	exitCode := m.Run()

	err = client.Database("backend-challenge-test").Drop(context.TODO())

	if err != nil {
		panic(err)
	}

	err = client.Disconnect(context.TODO())

	if err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}

func TestCreateCategory_Valid(t *testing.T) {
	categoryRepo := repository.NewCategoryRepository(testDB, "backend-challenge-test", "categories")

	category := &entities.Category{
		ID:            primitive.NewObjectID(),
		Name:          "Electronics",
		Subcategories: []string{"Phones", "Laptops"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Insert the category into the repository

	err := categoryRepo.Create(category)

	assert.NoError(t, err, "Expected no error when creating category")

	// Fetch the category from the repository

	fetchedCategory, err := categoryRepo.GetByID(category.ID.Hex())
	assert.NoError(t, err, "Expected no error when fetching category")
	assert.Equal(t, category.Name, fetchedCategory.Name, "category name should match")
	assert.Equal(t, category.Subcategories, fetchedCategory.Subcategories, "category subcategories should match")
}

func TestCreateCategory_MissingName(t *testing.T) {
	categoryRepo := repository.NewCategoryRepository(testDB, "backend-challenge-test", "categories")

	category := &entities.Category{
		ID:            primitive.NewObjectID(),
		Subcategories: []string{"Phones", "Laptops"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := categoryRepo.Create(category)

	assert.Error(t, err, "expected error when creating category with missing name")
}

func TestCreateCategory_DuplicateID(t *testing.T) {

	categoryRepo := repository.NewCategoryRepository(testDB, "backend-challenge-test", "categories")

	category := &entities.Category{
		ID:            primitive.NewObjectID(),
		Name:          "Clothing",
		Subcategories: []string{"Men", "Women"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := categoryRepo.Create(category)

	assert.NoError(t, err, "expected no error when creating first category")

	duplicateCategory := &entities.Category{
		ID:            category.ID,
		Name:          "Footwear",
		Subcategories: []string{"Shoes", "Sandals"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = categoryRepo.Create(duplicateCategory)

	assert.Error(t, err, "expected error when creating duplicate category")
}

func TestCreateCategory_EmptySubcategories(t *testing.T) {
	categoryRepo := repository.NewCategoryRepository(testDB, "backend-challenge-test", "categories")

	category := &entities.Category{
		ID:            primitive.NewObjectID(),
		Name:          "Accesories",
		Subcategories: []string{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := categoryRepo.Create(category)

	assert.NoError(t, err, "expected no error when creating category with empty subcategories")

	fetchedCategory, err := categoryRepo.GetByID(category.ID.Hex())
	assert.NoError(t, err, "expected no error when fetching category")

	assert.Equal(t, category.Name, fetchedCategory.Name, "category name should match")
	assert.Empty(t, fetchedCategory.Subcategories, "category subcategories should be empty")

}
