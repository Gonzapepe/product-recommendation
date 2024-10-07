package handlers

import (
	"backend-challenge/internal/application/services"
	"backend-challenge/internal/domain/entities"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.categoryService.GetCategoryByID(id)

	if err != nil {
		HandleError(c, http.StatusNotFound, err)
		return
	}
	
	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context)  {
	categories, err := h.categoryService.GetAllCategories()

	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category entities.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	validationErrors := entities.ValidateStruct(&category)

	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationErrors": validationErrors})
		return
	}

	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	err := h.categoryService.CreateCategory(&category); if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}
	
	c.JSON(http.StatusCreated, category)
	
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var category entities.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	validationErrors := entities.ValidateStruct(&category)

	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationErrors": validationErrors})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	category.ID = objectId

	err = h.categoryService.UpdateCategory(&category); if err != nil {
		HandleError(c, http.StatusInternalServerError, err)

		log.Printf("ERROR: %v", err)
		
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context)  {
	id := c.Param("id")	

	if err := h.categoryService.DeleteCategory(id); err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}