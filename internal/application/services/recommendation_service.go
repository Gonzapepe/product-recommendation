package services

import (
	"backend-challenge/internal/domain/entities"
	"math"
	"sort"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

type Recommendation struct {
	Product         *entities.Product `json:"product"`
	SimilarityScore float64           `json:"similarity_score"`
}

type RecommendationService struct{}

func NewRecommendationService() *RecommendationService {
	return &RecommendationService{}
}

// RecommendSimilarProducts recommends similar products based on a target product
func (s *RecommendationService) RecommendSimilarProducts(targetProduct entities.Product, allProducts []*entities.Product) []*Recommendation {
	targetVector := s.ExtractFeatureVector(targetProduct)

	recommendations := []*Recommendation{}

	// Calculate similarity for all products
	for _, product := range allProducts {
		if product.ID != targetProduct.ID { // Exclude the target product itself
			productVector := s.ExtractFeatureVector(*product)
			similarity := s.CosineSimilarity(targetVector, productVector)

			recommendations = append(recommendations, &Recommendation{
				Product:         product,
				SimilarityScore: similarity,
			})
		}
	}

	// Sort by similarity score
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].SimilarityScore > recommendations[j].SimilarityScore
	})

	// Return top 5 recommendations.
	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
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
	nameTokens := s.tokenizeLocalizedString(product.Name.LocalizedString)
	descriptionTokens := s.tokenizeLocalizedString(product.Description.LocalizedString)

	// Combine tokens
	tokens := append(nameTokens, descriptionTokens...)

	

	for _, token := range tokens {
		features["word_"+token] += 1.0
	}

	return features
}

func (s *RecommendationService) tokenizeLocalizedString(localized entities.LocalizedString) []string {
	var tokens []string

	if localized.En != nil {
		tokens = append(tokens, tokenize(*localized.En, "en")...)
	}

	if localized.Es != nil {
		tokens = append(tokens, tokenize(*localized.Es, "es")...)
	}

	if localized.Pt != nil {
		tokens = append(tokens, tokenize(*localized.Pt, "pt")...)
	}

	return tokens
}

// normalize scales a value to a range of 0 to 1
func normalize(value float64) float64 {
	if value == 0 {
		return 0
	}

	return 1.0 / (1.0 + math.Exp(-value)) // Sigmoid normalization
}

// Stopwords for different languages
var stopwords = map[string]map[string]bool{
	"en": {
		"the": true, "is": true, "and": true, "a": true, "of": true, "to": true,
	},
	"es": {
		"el": true, "es": true, "y": true, "un": true, "de": true, "para": true,
	},
	"pt": {
		"o": true, "é": true, "e": true, "um": true, "de": true, "para": true,
	},
}


// tokenize processes text by removing stopwords, handling punctuation, and applying stemming
func tokenize(text, lang string) []string {
	// Default to English stopwords if no language is specified
	langStopwords, exists := stopwords[lang]
	if !exists {
		langStopwords = stopwords["en"]
	}

	// Split text into words and remove punctuation
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) // Remove non-alphanumeric
	})

	// Filter out stopwords and apply stemming
	var tokens []string
	for _, word := range words {
		if !langStopwords[word] { // Skip stopwords
			stemmed, err := snowball.Stem(word, lang, true) // Apply stemming for the specified language
			if err == nil {
				tokens = append(tokens, stemmed)
			} else {
				tokens = append(tokens, word) // Fallback to original word
			}
		}
	}

	return tokens
}