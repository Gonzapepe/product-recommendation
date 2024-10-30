package tests

// import (
// 	"backend-challenge/internal/adapters/persistence/repository"
// 	"backend-challenge/internal/domain/entities"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func stringPtr(s string) *string {
//     return &s
// }

// func TestProductCreation(t *testing.T) {

// 	productRepo := repository.NewProductRepository(testDB, "backend-challenge-test", "products")

// 	product := &entities.Product{
// 		ID:          "test-id",
// 		StoreID:     "test-store-id",
// 		Categories: []entities.CategoryForProduct{
// 			{
// 				ID: "test-category-id",
// 				Name: entities.CategoryName{
// 					LocalizedString: entities.LocalizedString{
// 						En: stringPtr("Test Category"),
// 					},
// 				},
// 				Subcategories: []string{"subcategory1", "subcategory2"},
// 			},
// 		},
// 		Description: entities.Description{
// 			LocalizedString: entities.LocalizedString{
// 				En: stringPtr("Test Description"),
// 			},
// 		},
// 		Images: []entities.Image{
// 			{
// 				ID:       1,
// 				Src:      "http://example.com/image.png",
// 				Position: 1,
// 				Alt:      []entities.Alt{{LocalizedString: entities.LocalizedString{En: stringPtr("Test Alt")}}},
// 			},
// 		},
// 		Name: entities.Name{
// 			LocalizedString: entities.LocalizedString{
// 				En: stringPtr("Test Product"),
// 			},
// 		},
// 		Published:   true,
// 		Urls:        entities.Urls{CanonicalURL: "http://example.com/product", VideoURL: stringPtr("http://example.com/video.mp4")},
// 		Variants:    []entities.Variant{{ID: "variant-id", Value: "test-value", Stock: 100, Price: 29.99}},
// 		SoldCount:   50,
// 		ClickCount:  10,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}

// 	err := productRepo.Create(product)

// 	assert.NoError(t, err)

// 	savedProduct, err := productRepo.GetByID("test-id")

// 	// Assertions for all fields
// 	assert.NoError(t, err)
// 	assert.NotNil(t, savedProduct)
// 	assert.Equal(t, product.ID, savedProduct.ID)
// 	assert.Equal(t, product.StoreID, savedProduct.StoreID)
// 	assert.Equal(t, product.Categories[0].ID, savedProduct.Categories[0].ID)
// 	assert.Equal(t, product.Description.LocalizedString.En, savedProduct.Description.LocalizedString.En)
// 	assert.Equal(t, product.Images[0].Src, savedProduct.Images[0].Src)
// 	assert.Equal(t, product.Name.LocalizedString.En, savedProduct.Name.LocalizedString.En)
// 	assert.Equal(t, product.Published, savedProduct.Published)
// 	assert.Equal(t, product.Urls.CanonicalURL, savedProduct.Urls.CanonicalURL)
// 	assert.Equal(t, product.Variants[0].ID, savedProduct.Variants[0].ID)
// 	assert.Equal(t, product.SoldCount, savedProduct.SoldCount)
// 	assert.Equal(t, product.ClickCount, savedProduct.ClickCount)

// }

// func TestProductDeletion(t *testing.T) {
// 	productRepo := repository.NewProductRepository(testDB, "backend-challenge-test", "products")

// 	product := &entities.Product{
// 		ID:          "test-id",
// 		StoreID:     "test-store-id",
// 		Categories: []entities.CategoryForProduct{
// 			{
// 				ID: "test-category-id",
// 				Name: entities.CategoryName{
// 					LocalizedString: entities.LocalizedString{
// 						En: stringPtr("Test Category"),
// 					},
// 				},
// 				Subcategories: []string{"subcategory1", "subcategory2"},
// 			},
// 		},
// 		Description: entities.Description{
// 			LocalizedString: entities.LocalizedString{
// 				En: stringPtr("Test Description"),
// 			},
// 		},
// 		Images: []entities.Image{
// 			{
// 				ID:       1,
// 				Src:      "http://example.com/image.png",
// 				Position: 1,
// 				Alt:      []entities.Alt{{LocalizedString: entities.LocalizedString{En: stringPtr("Test Alt")}}},
// 			},
// 		},
// 		Name: entities.Name{
// 			LocalizedString: entities.LocalizedString{
// 				En: stringPtr("Test Product"),
// 			},
// 		},
// 		Published:   true,
// 		Urls:        entities.Urls{CanonicalURL: "http://example.com/product", VideoURL: stringPtr("http://example.com/video.mp4")},
// 		Variants:    []entities.Variant{{ID: "variant-id", Value: "test-value", Stock: 100, Price: 29.99}},
// 		SoldCount:   50,
// 		ClickCount:  10,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}

// 	err := productRepo.Create(product)

// 	assert.NoError(t, err)

// 	err = productRepo.Delete(product.ID)

// 	assert.NoError(t, err)
// }

// func TestProductUpdate(t *testing.T) {
// 	productRepo := repository.NewProductRepository(testDB, "backend-challenge-test", "products")

// 	initialProduct := &entities.Product{
// 		ID:          "test-id",
// 		StoreID:     "test-store-id",
// 		Categories: []entities.CategoryForProduct{
// 			{
// 				ID: "test-category-id",
// 				Name: entities.CategoryName{
// 					LocalizedString: entities.LocalizedString{
// 						En: stringPtr("Test Category"),
// 					},
// 				},
// 				Subcategories: []string{"subcategory1", "subcategory2"},
// 			},
// 		},
// 		Description: entities.Description{
// 			LocalizedString: entities.LocalizedString{
// 				En: stringPtr("Test Description"),
// 			},
// 		},
// 		Images: []entities.Image{
// 			{
// 				ID:       1,
// 				Src:      "http://example.com/image.png",
// 				Position: 1,
// 				Alt:      []entities.Alt{{LocalizedString: entities.LocalizedString{En: stringPtr("Test Alt")}}},
// 			},
// 		},
// 		Name: entities.Name{
// 			LocalizedString: entities.LocalizedString{
// 				En: stringPtr("Test Product"),
// 			},
// 		},
// 		Published:   true,
// 		Urls:        entities.Urls{CanonicalURL: "http://example.com/product", VideoURL: stringPtr("http://example.com/video.mp4")},
// 		Variants:    []entities.Variant{{ID: "variant-id", Value: "test-value", Stock: 100, Price: 29.99}},
// 		SoldCount:   50,
// 		ClickCount:  10,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}

// 	err := productRepo.Create(initialProduct)
// 	require.NoError(t, err)

// 	updateProduct := &entities.Product{
// 		ID:          initialProduct.ID,
// 		Name:        entities.Name{
// 			LocalizedString: entities.LocalizedString{
// 				En: stringPtr("Updated Test Product"),
// 			},
// 		},
// 		Description: entities.Description{
// 			LocalizedString: entities.LocalizedString{
// 				En: stringPtr("Updated Test Description"),
// 			},
// 		},
// 		SoldCount: 75,
// 	}

// 	err = productRepo.Update(updateProduct)
// 	assert.NoError(t, err)

// 	// Verify the product was updated
// 	var updatedProduct *entities.Product

// 	updatedProduct, err = productRepo.GetByID(initialProduct.ID)
// 	require.NoError(t, err)

// 	assert.Equal(t, updateProduct.Name.LocalizedString.En, updatedProduct.Name.LocalizedString.En)
// 	assert.Equal(t, updateProduct.Description.LocalizedString.En, updatedProduct.Description.LocalizedString.En)
// 	assert.Equal(t, updateProduct.SoldCount, updatedProduct.SoldCount)
// }