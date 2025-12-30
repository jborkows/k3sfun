package views

import (
	"shopping/internal/domain/products"
	"shopping/internal/domain/shoppinglist"
	"shopping/internal/infrastructure/oidc"
)

// ProductGroup represents a group of products for display purposes.
type ProductGroup struct {
	Name     string
	Products []products.Product
}

// ShoppingItemGroup represents a group of shopping list items for display purposes.
type ShoppingItemGroup struct {
	Name  string
	Items []shoppinglist.Item
}

type BaseData struct {
	Title         string
	User          *oidc.User
	HTMXSrc       string
	StaticVersion string
	IsAdmin       bool
}

// HTMXSSESrc returns the SSE extension URL derived from the main HTMX source.
// For example, if HTMXSrc is "https://unpkg.com/htmx.org@1.9.12",
// this returns "https://unpkg.com/htmx.org@1.9.12/dist/ext/sse.js".
func (b BaseData) HTMXSSESrc() string {
	if b.HTMXSrc == "" {
		return "https://unpkg.com/htmx.org@1.9.12/dist/ext/sse.js"
	}
	return b.HTMXSrc + "/dist/ext/sse.js"
}

type ProductsPageData struct {
	Base             BaseData
	Groups           []products.Group
	Products         []products.Product
	OnlyMissing      bool
	NameQuery        string
	SelectedGroupIDs []products.GroupID
	Page             int64
	TotalPages       int64
	Total            int64
}

type ProductsListData struct {
	Groups           []products.Group
	Products         []products.Product
	OnlyMissing      bool
	NameQuery        string
	SelectedGroupIDs []products.GroupID
	Page             int64
	TotalPages       int64
	Total            int64
}

type AdminPageData struct {
	Base BaseData
}

type ShoppingListData struct {
	Items []shoppinglist.Item
	Units []products.Unit
}

type ShoppingListPageData struct {
	Base  BaseData
	Items []shoppinglist.Item
	Units []products.Unit
}

type ProductsNewPageData struct {
	Base   BaseData
	Groups []products.Group
}
