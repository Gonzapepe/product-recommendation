package services

import (
	"backend-challenge/internal/domain/entities"
	"math"
	"strings"
)

type RecommendationService struct {}

func NewRecommendationService() *RecommendationService {
	return &RecommendationService{}
}

// ExtractFeatureVector creates a feature vector for a product
func (s *RecommendationService) ExtractFeatureVector(product entities.Product) map[string]float64 {
	// Initialize feature map
	features := make(map[string]float64)

	// Add category features (one-hot encoding)
	for _, category := range product.Categories {
		features["category_"+strings.ToLower(category)] = 1.0
	}

	// Add popularity metrics 
	features["click_count"] = normalize(float64(product.ClickCount))
	features["sold_count"] = normalize(float64(product.SoldCount))

	// Add textual features (simplified; extend with NLP later)
	words := tokenize(s.getLocalizedStringText(product.Name.LocalizedString) + " " +s.getLocalizedStringText(product.Description.LocalizedString)) 

	for _, word := range words {
		features["word_"+strings.ToLower(word)] += 1.0
	}

	return features
}

func (s *RecommendationService) getLocalizedStringText(localized entities.LocalizedString) string {
	var texts []string

	if localized.En != nil {
		texts = append(texts, *localized.En)
	}

	if localized.Es != nil {
		texts = append(texts, *localized.Es)
	}

	if localized.Pt != nil {
		texts = append(texts, *localized.Pt)
	}

	return strings.Join(texts, " ")
}

// normalize scales a value to a range of 0 to 1
func normalize(value float64) float64 {
	if value == 0 {
		return 0
	}

	return 1.0 / (1.0 + math.Exp(-value)) // Sigmoid normalization
}

func tokenize(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

