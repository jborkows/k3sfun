package shoppinglist

import (
	"time"

	"shopping/internal/domain/products"
)

type ItemID int64

type Item struct {
	ID           ItemID
	Name         string
	ProductID    *products.ProductID
	IconKey      string
	GroupName    string
	GroupOrder   int64
	Quantity     products.Quantity
	Unit         products.Unit
	UnitSingular string
	UnitPlural   string
	Done         bool
	IntegerOnly  bool
	CreatedAt    time.Time
}
