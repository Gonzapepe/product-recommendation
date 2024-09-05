package entities

import "time"

type BrainBoundary struct {
	PriceLimit          *float64
	CategoryRestriction []string
	OnlyInStock         bool
	RecentlyAdded       bool
}

type BrainRule struct {
	PrioritizeCategories []string
	SortBy               string
}

type BrainMetadata struct {
	ProcessingTime        time.Duration
	ProductsConsidered    int
	ProductsReturned      int
	BoundariesApplied     []BrainBoundary
	RulesApplied          []BrainRule
	FilteredProductsCount int
}
