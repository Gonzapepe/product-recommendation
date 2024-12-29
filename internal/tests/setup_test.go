// test_setup.go
package tests

import (
	"backend-challenge/internal/adapters/persistence/repository"
	"backend-challenge/internal/adapters/web/handlers"
	"backend-challenge/internal/application/services"
	"backend-challenge/internal/domain/repositories"
	"context"
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testDB          *mongo.Database
	testClient      *mongo.Client
	categoryRepo    repositories.CategoryRepository
	productRepo     repositories.ProductRepository
	recommendationService *services.RecommendationService
	categoryService services.CategoryService
	productService  services.ProductService
	categoryHandler *handlers.CategoryHandler
	productHandler  *handlers.ProductHandler
)

// TestMain is the main entry point for tests in this package
func TestMain(m *testing.M) {

	fmt.Println("TestMain is running")

	// Setup
	mongoURI := os.Getenv("TEST_MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	testClient = client // set the client globally
	testDB = client.Database("backend-challenge-test")

	// Initialize repositories and services
	categoryRepo = repository.NewCategoryRepository(testClient, "backend-challenge-test", "categories")
	productRepo = repository.NewProductRepository(testClient, "backend-challenge-test", "products")
	recommendationService = services.NewRecommendationService()
	categoryService = services.NewCategoryService(categoryRepo)
	productService = services.NewProductService(productRepo, recommendationService)

	// Initialize handlers
	categoryHandler = handlers.NewCategoryHandler(categoryService)
	productHandler = handlers.NewProductHandler(productService /*brainService*/)

	// Run tests
	exitCode := m.Run()

	// Cleanup
	err = testDB.Drop(context.TODO())
	if err != nil {
		panic(err)
	}
	err = client.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}

// setupTest is a helper function that can be used by individual test files
func setupTest(t *testing.T) func() {
	// Return a cleanup function
	return func() {
		// Add any cleanup code here
		collections, err := testDB.ListCollectionNames(context.TODO(), map[string]interface{}{})
		if err != nil {
			t.Fatal(err)
		}

		for _, collection := range collections {
			err := testDB.Collection(collection).Drop(context.TODO())
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

// GetTestClient provides access to the test client
func GetTestClient() *mongo.Client {
	if testClient == nil {
		panic("Test client not initialized")
	}
	return testClient
}

func GetCategoryHandler() *handlers.CategoryHandler {
	if categoryHandler == nil {
		panic("Category repository not initialized")
	}
	return categoryHandler
}

func GetProductRepo() repositories.ProductRepository {
	if productRepo == nil {
		panic("Category repository not initialized")
	}
	return productRepo
}

func GetProductHandler() *handlers.ProductHandler {
	if productHandler == nil {
		panic("Category repository not initialized")
	}
	return productHandler
}

func GetCategoryRepo() repositories.CategoryRepository {
	if categoryRepo == nil {
		panic("Category repository not initialized")
	}
	return categoryRepo
}

func GetRecommendationService() *services.RecommendationService {
	if recommendationService == nil {
		panic("Category repository not initialized")
	}
	return recommendationService
}