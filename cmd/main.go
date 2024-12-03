package main

import (
	"context"
	"log"

	"backend-challenge/internal/adapters/persistence/repository"
	"backend-challenge/internal/adapters/web/handlers"
	"backend-challenge/internal/application/services"
	"backend-challenge/internal/infrastructure/db"

	"github.com/gin-gonic/gin"
)

var (
	productService services.ProductService
	categoryService services.CategoryService
	// brainService   *services.BrainService
)

func main() {
	mongoURI := "mongodb://localhost:27017"

	mongoClient, err := db.ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	productRepo := repository.NewProductRepository(mongoClient, "backend-challenge", "products")

	recommendationService := services.NewRecommendationService()

	productService = services.NewProductService(productRepo, recommendationService)

	categoryRepo := repository.NewCategoryRepository(mongoClient, "backend-challenge", "categories")
	
	categoryService  = services.NewCategoryService(categoryRepo)
	
	// Amount of product recommendations
	// brainService = services.NewBrainService(15)



	InitRoutes()
}

func InitRoutes() {
	router := gin.Default()

	v1 := router.Group("/v1")

	productHandler := handlers.NewProductHandler(productService, /*brainService*/)

	v1.POST("/products", productHandler.CreateProduct)
	v1.GET("/products", productHandler.GetAllProducts)
	v1.GET("/products/:id", productHandler.GetProductByID)
	v1.PUT("/products/:id", productHandler.UpdateProduct)
	v1.DELETE("products/:id", productHandler.DeleteProduct)

	categoryHandler := handlers.NewCategoryHandler(categoryService)

	v1.POST("/categories", categoryHandler.CreateCategory)
	v1.GET("/categories", categoryHandler.GetAllCategories)
	v1.GET("/categories/:id", categoryHandler.GetCategoryByID)
	v1.PUT("/categories/:id", categoryHandler.UpdateCategory)
	v1.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	router.Run(":8080")
}
