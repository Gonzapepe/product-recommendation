package tests

import (
	"backend-challenge/internal/domain/entities"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateCategory_Valid(t *testing.T) {

	// Setup test and defer cleanup
	cleanup := setupTest(t)
	defer cleanup()

	// Create the router
	router := gin.Default()

	categoryRepo := GetCategoryRepo()
	categoryHandler := GetCategoryHandler()

	router.POST("/categories", categoryHandler.CreateCategory)

	category := &entities.Category{
		ID:            primitive.NewObjectID(),
		Name:          "Electronics",
		Subcategories: []string{"Phones", "Laptops"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Insert the category into the repository

	categoryJSON, _ := json.Marshal(category)

	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(categoryJSON))

	req.Header.Set("Content-Type", "application/json")

	// Execute the request with httptest
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code, "Expected status code 201 for successful category creation")

	var responseMap map[string]interface{}

	err := json.Unmarshal(resp.Body.Bytes(), &responseMap)

	fmt.Println(responseMap)
	assert.NoError(t, err, "Expected no error when unmarshalling response")
	id, ok := responseMap["id"].(string)
	assert.True(t, ok, "Expected id to be a string")

	// Fetch the category from the repository

	fetchedCategory, err := categoryRepo.GetByID(id)
	assert.NoError(t, err, "Expected no error when fetching category")
	assert.Equal(t, category.Name, fetchedCategory.Name, "category name should match")
	assert.Equal(t, category.Subcategories, fetchedCategory.Subcategories, "category subcategories should match")
}

func TestCreateCategory_MissingName(t *testing.T) {

	// Setup test and defer cleanup
	cleanup := setupTest(t)
	defer cleanup()

	categoryHandler := GetCategoryHandler()

	router := gin.Default()

	router.POST("/categories", categoryHandler.CreateCategory)

	category := &entities.Category{
		ID:            primitive.NewObjectID(),
		Subcategories: []string{"Phones", "Laptops"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	categoryJSON, _ := json.Marshal(category)

	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(categoryJSON))

	req.Header.Set("Content-Type", "application/json")

	// Execute the request with httptest
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Expected status code 400 for missing name")
}

func TestCreateCategory_EmptySubcategories(t *testing.T) {

	// Setup test and defer cleanup
	cleanup := setupTest(t)
	defer cleanup()

	categoryHandler := GetCategoryHandler()

	router := gin.Default()

	router.POST("/categories", categoryHandler.CreateCategory)

	category := &entities.Category{
		ID:            primitive.NewObjectID(),
		Name:          "Accesories",
		Subcategories: []string{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	categoryJSON, _ := json.Marshal(category)

	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(categoryJSON))

	req.Header.Set("Content-Type", "application/json")

	// Execute the request with httptest
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code, "Expected status code 201 for empty subcategories")

	var responseMap map[string]interface{}

	_ = json.Unmarshal(resp.Body.Bytes(), &responseMap)

	id := responseMap["id"].(string)

	fetchedCategory, err := categoryRepo.GetByID(id)
	assert.NoError(t, err, "expected no error when fetching category")

	assert.Equal(t, category.Name, fetchedCategory.Name, "category name should match")
	assert.Empty(t, fetchedCategory.Subcategories, "category subcategories should be empty")

}

func TestCategoryCreation(t *testing.T) {
	// Setup test and defer cleanup
	cleanup := setupTest(t)
	defer cleanup()

	categoryRepo := GetCategoryRepo()

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


	savedCategory, err := categoryRepo.GetByID(category.ID.Hex())

	assert.NoError(t, err, "Expected no error when fetching category")

	assert.Equal(t, category.Name, savedCategory.Name, "category name should match")
	assert.Equal(t, category.Subcategories, savedCategory.Subcategories, "category subcategories should match")

}

func TestGetAllCategories(t *testing.T) {
	// Setup test and defer cleanup
	cleanup := setupTest(t)
	defer cleanup()

	categoryRepo := GetCategoryRepo()

	categoryA := &entities.Category{
		ID:            primitive.NewObjectID(),
		Name:          "Electronics",
		Subcategories: []string{"Phones", "Laptops"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	categoryB := &entities.Category{
		ID:            primitive.NewObjectID(),
		Name:          "Electronics B",
		Subcategories: []string{"TVs", "Soundbars"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	categoryC := &entities.Category{
		ID:            primitive.NewObjectID(),
		Name:          "Electronics C",
		Subcategories: []string{"Gaming Consoles", "Speakers"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Insert the categories into the repository
	err := categoryRepo.Create(categoryA)
	assert.NoError(t, err, "Expected no error when creating category")

	err = categoryRepo.Create(categoryB)
	assert.NoError(t, err, "Expected no error when creating category")

	err = categoryRepo.Create(categoryC)
	assert.NoError(t, err, "Expected no error when creating category")

	// Get all categories from the repository
	categories, err := categoryRepo.GetAll()

	fmt.Printf("returned categories: %+v\n", categories)

	assert.NoError(t, err, "Expected no error when fetching categories")
	assert.Equal(t, 3, len(categories), "Expected 3 categories to be returned")
	assert.Condition(t, func() bool {
		for _, category := range categories {
			if category.Name == categoryA.Name &&
				category.Subcategories[0] == categoryA.Subcategories[0] {
				return true
			}
		}
		return false
	}, "Expected category A to be returned")

}

func TestUpdateCategory(t *testing.T) {
	// Setup test and defer cleanup
	cleanup := setupTest(t)	
	defer cleanup()

	categoryRepo := GetCategoryRepo()

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

	// Update the category
	category.Name = "Electronics B"
	err = categoryRepo.Update(category)
	assert.NoError(t, err, "Expected no error when updating category")

	// Fetch the updated category from the repository
	fetchedCategory, err := categoryRepo.GetByID(category.ID.Hex())
	assert.NoError(t, err, "Expected no error when fetching category")

	assert.Equal(t, category.Name, fetchedCategory.Name, "category name should match")
}

func TestDeleteCategory(t *testing.T) {
	// Setup test and defer cleanup
	cleanup := setupTest(t)
	defer cleanup()

	categoryRepo := GetCategoryRepo()

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

	// Delete the category	
	err = categoryRepo.Delete(category.ID.Hex())
	assert.NoError(t, err, "Expected no error when deleting category")

	// Fetch the deleted category from the repository
	_, err = categoryRepo.GetByID(category.ID.Hex())
	assert.Error(t, err, "Expected error when fetching deleted category")
}