package tests

import (
	"backend-challenge/internal/adapters/persistence/repository"
	"backend-challenge/internal/application/services"
	"backend-challenge/internal/domain/entities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func TestProductCreation(t *testing.T) {

	productRepo := repository.NewProductRepository(testClient, "backend-challenge-test", "products")

	product := &entities.Product{
		ID:          primitive.NewObjectID(),
		StoreID:     "test-store-id",
		Categories: []string{
			"Test Category",
		},
		Description: entities.Description{
			LocalizedString: entities.LocalizedString{
				En: ptr("Test Description"),
			},
		},
		Images: []entities.Image{
			{
				ID:       1,
				Src:      "http://example.com/image.png",
				Position: 1,
				Alt:      []entities.Alt{{LocalizedString: entities.LocalizedString{En: ptr("Test Alt")}}},
			},
		},
		Name: entities.Name{
			LocalizedString: entities.LocalizedString{
				En: ptr("Test Product"),
			},
		},
		Published:   true,
		Urls:        entities.Urls{CanonicalURL: "http://example.com/product", VideoURL: ptr("http://example.com/video.mp4")},
		Variants:    []entities.Variant{{ID: "variant-id", Value: "test-value", Stock: 100, Price: 29.99}},
		SoldCount:   50,
		ClickCount:  10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := productRepo.Create(product)

	assert.NoError(t, err)

	savedProduct, err := productRepo.GetByID(product.ID.Hex())

	// Assertions for all fields
	assert.NoError(t, err)
	assert.NotNil(t, savedProduct)
	assert.Equal(t, product.ID, savedProduct.ID)
	assert.Equal(t, product.StoreID, savedProduct.StoreID)
	assert.Equal(t, product.Categories[0], savedProduct.Categories[0])
	assert.Equal(t, product.Description.LocalizedString.En, savedProduct.Description.LocalizedString.En)
	assert.Equal(t, product.Images[0].Src, savedProduct.Images[0].Src)
	assert.Equal(t, product.Name.LocalizedString.En, savedProduct.Name.LocalizedString.En)
	assert.Equal(t, product.Published, savedProduct.Published)
	assert.Equal(t, product.Urls.CanonicalURL, savedProduct.Urls.CanonicalURL)
	assert.Equal(t, product.Variants[0].ID, savedProduct.Variants[0].ID)
	assert.Equal(t, product.SoldCount, savedProduct.SoldCount)
	assert.Equal(t, product.ClickCount, savedProduct.ClickCount)

}

func TestProductDeletion(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	productRepo := GetProductRepo()

	product := &entities.Product{
		ID:          primitive.NewObjectID(),
		StoreID:     "test-store-id",
		Categories: []string{
			"Test Category",
		},
		Description: entities.Description{
			LocalizedString: entities.LocalizedString{
				En: ptr("Test Description"),
			},
		},
		Images: []entities.Image{
			{
				ID:       1,
				Src:      "http://example.com/image.png",
				Position: 1,
				Alt:      []entities.Alt{{LocalizedString: entities.LocalizedString{En: ptr("Test Alt")}}},
			},
		},
		Name: entities.Name{
			LocalizedString: entities.LocalizedString{
				En: ptr("Test Product"),
			},
		},
		Published:   true,
		Urls:        entities.Urls{CanonicalURL: "http://example.com/product", VideoURL: ptr("http://example.com/video.mp4")},
		Variants:    []entities.Variant{{ID: "variant-id", Value: "test-value", Stock: 100, Price: 29.99}},
		SoldCount:   50,
		ClickCount:  10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := productRepo.Create(product)

	assert.NoError(t, err)

	err = productRepo.Delete(product.ID.Hex())

	assert.NoError(t, err)
}

func TestProductUpdate(t *testing.T) {
	productRepo := repository.NewProductRepository(testClient, "backend-challenge-test", "products")

	initialProduct := &entities.Product{
		ID:          primitive.NewObjectID(),
		StoreID:     "test-store-id",
		Categories: []string{"Test Category"},
		Description: entities.Description{
			LocalizedString: entities.LocalizedString{
				En: ptr("Test Description"),
			},
		},
		Images: []entities.Image{
			{
				ID:       1,
				Src:      "http://example.com/image.png",
				Position: 1,
				Alt:      []entities.Alt{{LocalizedString: entities.LocalizedString{En: ptr("Test Alt")}}},
			},
		},
		Name: entities.Name{
			LocalizedString: entities.LocalizedString{
				En: ptr("Test Product"),
			},
		},
		Published:   true,
		Urls:        entities.Urls{CanonicalURL: "http://example.com/product", VideoURL: ptr("http://example.com/video.mp4")},
		Variants:    []entities.Variant{{ID: "variant-id", Value: "test-value", Stock: 100, Price: 29.99}},
		SoldCount:   50,
		ClickCount:  10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := productRepo.Create(initialProduct)
	require.NoError(t, err)

	updateProduct := &entities.Product{
		ID:          initialProduct.ID,
		Name:        entities.Name{
			LocalizedString: entities.LocalizedString{
				En: ptr("Updated Test Product"),
			},
		},
		Description: entities.Description{
			LocalizedString: entities.LocalizedString{
				En: ptr("Updated Test Description"),
			},
		},
		SoldCount: 75,
	}

	err = productRepo.Update(updateProduct)
	assert.NoError(t, err)

	// Verify the product was updated
	var updatedProduct *entities.Product

	updatedProduct, err = productRepo.GetByID(initialProduct.ID.Hex())
	require.NoError(t, err)

	assert.Equal(t, updateProduct.Name.LocalizedString.En, updatedProduct.Name.LocalizedString.En)
	assert.Equal(t, updateProduct.Description.LocalizedString.En, updatedProduct.Description.LocalizedString.En)
	assert.Equal(t, updateProduct.SoldCount, updatedProduct.SoldCount)
}

func TestGetRecommendations(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()


	productRepo := GetProductRepo()

	recommendationService := services.NewRecommendationService()

	productService := services.NewProductService(productRepo, recommendationService)

	// Mock data
	productA := &entities.Product{ID: primitive.NewObjectID(), Categories: []string{"Electronics"}}
	productB := &entities.Product{ID: primitive.NewObjectID(), Categories: []string{"Electronics"}}
	productC := &entities.Product{ID: primitive.NewObjectID(), Categories: []string{"Home Appliances"}}

	err := productRepo.Create(productA)
	assert.NoError(t, err)

	err = productRepo.Create(productB)
	assert.NoError(t, err)

	err = productRepo.Create(productC)	
	assert.NoError(t, err)
	
	product, getIDErr :=productRepo.GetByID(productA.ID.Hex())

	assert.NoError(t, getIDErr)

	assert.Equal(t, productA, product)

	allResults, allResultsErr := productRepo.GetAll()

	assert.NoError(t, allResultsErr)

	assert.Equal(t, []*entities.Product{productA, productB, productC}, allResults)

	recommendations, err := productService.GetRecommendations(productA.ID.Hex())

	assert.NoError(t, err)
	assert.Equal(t, recommendations[0].Product.ID, productB.ID)

}
