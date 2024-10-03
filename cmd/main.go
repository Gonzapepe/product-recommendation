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
	brainService   *services.BrainService
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

	productService = services.NewProductService(productRepo)

	// Amount of product recommendations
	brainService = services.NewBrainService(15)

	InitRoutes()
}

func InitRoutes() {
	router := gin.Default()

	v1 := router.Group("/v1")

	productHandler := handlers.NewProductHandler(productService, *brainService)

	v1.POST("/products", productHandler.CreateProduct)
	v1.GET("/products", productHandler.GetAllProducts)
	v1.GET("/products/:id", productHandler.GetProductByID)
	v1.PUT("/products/:id", productHandler.UpdateProduct)
	v1.DELETE("products/:id", productHandler.DeleteProduct)

	router.Run(":8080")
}
