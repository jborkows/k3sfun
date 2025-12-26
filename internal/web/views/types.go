package views

import (
	"shopping/internal/domain/products"
	"shopping/internal/domain/shoppinglist"
	"shopping/internal/infrastructure/oidc"
)

type BaseData struct {
	Title         string
	User          *oidc.User
	HTMXSrc       string
	StaticVersion string
	IsAdmin       bool
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
}

type ShoppingListPageData struct {
	Base  BaseData
	Items []shoppinglist.Item
}
