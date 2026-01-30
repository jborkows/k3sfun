package products

import "context"

// Repository is the write-side persistence port used by Service.
// Infrastructure (e.g., SQLite) implements this interface.
type Repository interface {
	CreateGroup(ctx context.Context, name string) (GroupID, error)
	CreateProduct(ctx context.Context, p NewProduct) (ProductID, error)

	SetProductQuantity(ctx context.Context, productID ProductID, qty Quantity) error
	AddProductQuantity(ctx context.Context, productID ProductID, delta Quantity) error
	SetProductGroup(ctx context.Context, productID ProductID, groupID *GroupID) error
	SetProductUnit(ctx context.Context, productID ProductID, unit Unit) error

	ResolveIconKeyForName(ctx context.Context, name string) (string, bool, error)
}
