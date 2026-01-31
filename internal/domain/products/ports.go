package products

import "context"

// Queries are read-side operations (CQRS).
type Queries interface {
	ListGroups(ctx context.Context) ([]Group, error)
	ListProducts(ctx context.Context, filter ProductFilter) ([]Product, error)
	// ListProductsAll returns all products (no paging). Used by UI helpers
	// when rendering single product cards to avoid pagination limits.
	ListProductsAll(ctx context.Context) ([]Product, error)
	SuggestProductsByName(ctx context.Context, query string, limit int64) ([]Product, error)
	CountProducts(ctx context.Context, filter ProductFilter) (int64, error)
	ListUnits(ctx context.Context) ([]Unit, error)
}
