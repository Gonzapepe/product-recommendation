package handlers

import (
	"backend-challenge/internal/application/services"
	"backend-challenge/internal/domain/entities"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	productService services.ProductService
	// brainService	services.BrainService
}

func generateNewID() string {
	return uuid.New().String()
}

func NewProductHandler(productService services.ProductService, /*brainService services.BrainService */) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		// brainService: *services.NewBrainService(15),
	}
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	productId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	product, err := h.productService.GetProductByID(productId.Hex())

	if err != nil {
		HandleError(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {

	// _, limit, offset := getPaginationParams(c)

	// paginatedProducts, err := h.productService.GetPaginatedProducts(offset, limit)

	allProducts, err := h.productService.GetAllProducts()

	log.Printf("Error retrieving products: %v", err)

	if err != nil {
		log.Printf("Error retrieving products: %v", err)
		HandleError(c, http.StatusInternalServerError, err)
		return
	}


	// boundaries, rules := parseRecommendationParams(c)

	// recommendations, metadata := h.brainService.GenerateProductSuggestions(allProducts, paginatedProducts, boundaries, rules)

	c.JSON(http.StatusOK, gin.H{
		"products": allProducts,
		// "recommendations": recommendations,
		// "recommendaiton_metadata": metadata,
	})
}

func getPaginationParams(c *gin.Context) (page, limit, offset int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset = (page - 1) * limit
	return
}

func parseRecommendationParams(c *gin.Context) ([]entities.BrainBoundary, []entities.BrainRule) {
	priceLimit := parsePriceLimit(c)
	onlyInStock := c.DefaultQuery("onlyInStock", "false") == "true"
	categories := c.QueryArray("categories")
	prioritizeCategories := c.QueryArray("prioritizeCategories")
	sortBy := c.DefaultQuery("sortBy", "none")

	log.Printf("Categories: %v", categories)
	log.Printf("Prioritize categories: %v", prioritizeCategories)

	boundaries := []entities.BrainBoundary{
		{
			PriceLimit: priceLimit,
			OnlyInStock: onlyInStock,
			CategoryRestriction: categories,
		},
	}

	rules := []entities.BrainRule{
		{
			PrioritizeCategories: prioritizeCategories,
			SortBy: sortBy,
		},
	}

	return boundaries, rules
}

func parsePriceLimit(c *gin.Context) *float64 {
	priceLimitStr := c.Query("priceLimit")
	if priceLimitStr == "" {
		return nil
	}
	limit, err := strconv.ParseFloat(priceLimitStr, 64)
	if err != nil {
		return nil
	}

	return &limit
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product entities.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	validationErrors := entities.ValidateStruct(&product)

	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationErrors": validationErrors})
		return
	}

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	if err := h.productService.CreateProduct(&product); err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product entities.Product
	if err := c.ShouldBindBodyWithJSON(&product); err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	validationErrors := entities.ValidateStruct(&product)

	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationErrors": validationErrors})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}
	
	product.ID = objectId

	if err := h.productService.UpdateProduct(&product); err != nil {
		
		log.Printf("ERROR: %v", err)

		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, product)

}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.productService.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	
	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) GetRecommendations(c *gin.Context) {
	productID := c.Param("id")

	recommendations, err := h.productService.GetRecommendations(productID)

	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

 	c.JSON(http.StatusOK, recommendations)
}