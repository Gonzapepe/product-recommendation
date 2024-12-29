package tests

import (
	"backend-challenge/internal/domain/entities"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func TestRecommendSimilarProducts(t *testing.T) {
	// setup test and defer cleanup
	cleanup := setupTest(t)
	defer cleanup()

	recommendationService := GetRecommendationService()

	productA := entities.Product{

		ID: primitive.NewObjectID(),
		Categories: []string{"Electronics"},
		ClickCount: 100,
		SoldCount: 50,
		Name:        entities.Name{LocalizedString: entities.LocalizedString{En: ptr("Laptop A")}},
		Description: entities.Description{LocalizedString: entities.LocalizedString{En: ptr("A powerful laptop.")}},
	}

	productB := entities.Product{
		ID: primitive.NewObjectID(),
		Categories:  []string{"Electronics"},
		ClickCount: 80,
		SoldCount: 40,
		Name:        entities.Name{LocalizedString: entities.LocalizedString{En: ptr("Laptop B")}},
		Description: entities.Description{LocalizedString: entities.LocalizedString{En: ptr("A compact laptop.")}},
	}

	productC := entities.Product{
		ID:          primitive.NewObjectID(),
		Categories:  []string{"Home Appliances"},
		ClickCount:  20,
		SoldCount:   10,
		Name:        entities.Name{LocalizedString: entities.LocalizedString{En: ptr("Vacuum Cleaner")}},
		Description: entities.Description{LocalizedString: entities.LocalizedString{En: ptr("A high-efficiency vacuum cleaner.")}},
	}

	allProducts := []*entities.Product{&productA, &productB, &productC}

	// Get recommendations
	recommendations := recommendationService.RecommendSimilarProducts(productA, allProducts)

	// Assert results
	assert.Equal(t, 2, len(recommendations))
	assert.Equal(t, productB.ID, recommendations[0].Product.ID)
	assert.True(t, recommendations[0].SimilarityScore > recommendations[1].SimilarityScore)
}

func TestCosineSimilarity(t *testing.T) {
	recommendationService := GetRecommendationService()

	vec1 := map[string]float64{"a": 1.0, "b": 1.0, "c": 0.0}
	vec2 := map[string]float64{"a": 1.0, "b": 1.0, "c": 1.0}

	similarity := recommendationService.CosineSimilarity(vec1, vec2)

	assert.True(t, similarity > 0.8) // vec1 and vec2 overlap heavily so expect high similarity score
}

func ptr(s string) *string {
	return &s
}