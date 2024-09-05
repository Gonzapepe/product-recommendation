package services

import (
	"backend-challenge/internal/domain/entities"
	"sort"
	"strings"
	"time"
)

type BrainService struct {
	maxRecommendations int
}

func NewBrainService(maxRecommendations int) *BrainService {
	return &BrainService{
		maxRecommendations: maxRecommendations,
	}
}

func (b *BrainService) GenerateProductSuggestions(products []*entities.Product, paginatedProducts []*entities.Product, boundaries []entities.BrainBoundary, rules []entities.BrainRule) ([]*entities.Product, entities.BrainMetadata) {
	startTime := time.Now()

	paginatedProductsIDs := make(map[string]bool)

	for _, p := range paginatedProducts {
		paginatedProductsIDs[p.ID] = true
	}

	candidateProducts := filter(products, func(p *entities.Product) bool {
		return !paginatedProductsIDs[p.ID]
	})

	filteredProducts := b.applyBoundaries(candidateProducts, boundaries)
	rankedProducts := b.applyRules(filteredProducts, rules)

	if len(rankedProducts) > b.maxRecommendations {
		rankedProducts = rankedProducts[:b.maxRecommendations]
	}

	metadata := entities.BrainMetadata{
		ProcessingTime:        time.Since(startTime),
		ProductsConsidered:    len(products),
		ProductsReturned:      len(rankedProducts),
		BoundariesApplied:     boundaries,
		RulesApplied:          rules,
		FilteredProductsCount: len(filteredProducts),
	}

	return rankedProducts, metadata
}

func (b *BrainService) applyBoundaries(products []*entities.Product, boundaries []entities.BrainBoundary) []*entities.Product {
	return filter(products, func(p *entities.Product) bool {
		for _, boundary := range boundaries {
			if !b.productMeetsBoundary(p, boundary) {
				return false
			}
		}
		return true
	})
}

func (b *BrainService) productMeetsBoundary(product *entities.Product, boundary entities.BrainBoundary) bool {
	if boundary.PriceLimit != nil && !b.meetsPriceLimit(product, *boundary.PriceLimit) {
		return false
	}
	if boundary.OnlyInStock && !b.hasStock(product) {
		return false
	}
	if len(boundary.CategoryRestriction) > 0 {
		hasMatchingCategory := false
		for _, p := range product.Categories {
			if p.Name.En != nil && contains(boundary.CategoryRestriction, *p.Name.En) ||
				p.Name.Es != nil && contains(boundary.CategoryRestriction, *p.Name.Es) ||
				p.Name.Pt != nil && contains(boundary.CategoryRestriction, *p.Name.Pt) {
				hasMatchingCategory = true
				break
			}
		}
		if !hasMatchingCategory {
			return false
		}
	}
	return true
}

func (b *BrainService) meetsPriceLimit(product *entities.Product, limit float64) bool {
	for _, variant := range product.Variants {
		if variant.Price <= limit {
			return true
		}
	}
	return false
}

func (b *BrainService) hasStock(product *entities.Product) bool {
	for _, variant := range product.Variants {
		if variant.Stock > 0 {
			return true
		}
	}
	return false
}

func (b *BrainService) applyRules(products []*entities.Product, rules []entities.BrainRule) []*entities.Product {

	prioritizedCategories := make([]string, 0)
	sortByPrice := false

	for _, rule := range rules {
		if len(rule.PrioritizeCategories) > 0 {
			prioritizedCategories = append(prioritizedCategories, rule.PrioritizeCategories...)
		}
		if rule.SortBy == "price" {
			sortByPrice = true
		}
	}

	if len(prioritizedCategories) > 0 {
		// split products into prioritized and non-prioritized
		prioritized := make([]*entities.Product, 0)
		others := make([]*entities.Product, 0)

		for _, p := range products {
			if productMatchesCategory(p, prioritizedCategories) {
				prioritized = append(prioritized, p)
			} else {
				others = append(others, p)
			}
		}

		if sortByPrice {
			prioritized = b.sortByPrice(prioritized)
			others = b.sortByPrice(others)

		}

		// combine prioritized and non-prioritized
		return append(prioritized, others...)
	} else if sortByPrice {
		// if only price sorting is needed
		return b.sortByPrice(products)
	}

	// if no rules are applied, return all products
	return products
}

func productMatchesCategory(product *entities.Product, categories []string) bool {
	for _, productCategory := range product.Categories {
        for _, category := range categories {
            if productCategory.Name.Es != nil && strings.EqualFold(*productCategory.Name.Es, category) {
                return true
            }
            if productCategory.Name.En != nil && strings.EqualFold(*productCategory.Name.En, category) {
                return true
            }
            if productCategory.Name.Pt != nil && strings.EqualFold(*productCategory.Name.Pt, category) {
                return true
            }
        }
    }
    return false
}

func (b *BrainService) sortByPrice(products []*entities.Product) []*entities.Product {
	sort.Slice(products, func(i, j int) bool {
		return getLowestVariantPrice(products[i].Variants) < getLowestVariantPrice(products[j].Variants)
	})
	return products
}


func getLowestVariantPrice(variants []entities.Variant) float64 {
	if len(variants) == 0 {
		return 0
	}
	lowestPrice := variants[0].Price
	for _, v := range variants {
		if v.Price < lowestPrice {
			lowestPrice = v.Price
		}
	}

	return lowestPrice
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if strings.EqualFold(s, needle) {
			return true
		}
	}
	return false
}

func filter(products []*entities.Product, predicate func(*entities.Product) bool) []*entities.Product {
	filtered := make([]*entities.Product, 0)
	for _, product := range products {
		if predicate(product) {
			filtered = append(filtered, product)
		}
	}

	return filtered
}
