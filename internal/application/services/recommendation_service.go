package services

import (
	"backend-challenge/internal/domain/entities"
	"math"
	"sort"
	"strings"
)

type RecommendationService struct {}

func NewRecommendationService() *RecommendationService {
	return &RecommendationService{}
}

// RecommendSimilarProducts recommends similar products based on a target product
func (s *RecommendationService) RecommendSimilarProducts (targetProduct entities.Product, allProducts []*entities.Product) []*entities.Product {
	targetVector := s.ExtractFeatureVector(targetProduct)

	similarities := make([]struct {
		Product *entities.Product
		Similarity float64
	}, len(allProducts))


	// Calculate similarity for all products
	for i, product := range allProducts {
		if product.ID != targetProduct.ID { // Exclude the target product itself
			productVector := s.ExtractFeatureVector(*product)
			similarities[i] = struct{Product *entities.Product; Similarity float64}{
				Product: product,
				Similarity: s.CosineSimilarity(targetVector, productVector),
			}
		}
	}

	// Sort by similarity score
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	// Return top 5 recommendations.
	var recommendations []*entities.Product
	for i := 0; i < len(similarities) && i < 5; i++ {
		recommendations = append(recommendations, similarities[i].Product)
	}

	return recommendations
}

// CosineSimilarity computes the cosine similarity between two feature vectors
func (s *RecommendationService) CosineSimilarity(vec1, vec2 map[string]float64) float64 {
	// Extracted info on cosine similarity:
	// https://en.wikipedia.org/wiki/Cosine_similarity

	var dotProduct, magnitudeA, magnitudeB float64

	for key, valueA := range vec1 {
		valueB, ok := vec2[key]
		if ok {
			dotProduct += valueA * valueB
		}

		magnitudeA += valueA * valueA
	}

	for _, valueB := range vec2 {
		magnitudeB += valueB * valueB
	}

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0 // Avoid division by zero
	}

	return dotProduct / (math.Sqrt(magnitudeA) * math.Sqrt(magnitudeB))
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

