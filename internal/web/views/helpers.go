package views

import (
	"net/url"
	"shopping/internal/domain/products"
	"shopping/internal/domain/shoppinglist"
	"strconv"
	"strings"
)

func boolToQuery(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func productsListQS(onlyMissing bool, nameQuery string, groupIDs []products.GroupID, page int64) string {
	values := url.Values{}
	if onlyMissing {
		values.Set("missing", "1")
	}
	nameQuery = strings.TrimSpace(nameQuery)
	if nameQuery != "" {
		values.Set("q", nameQuery)
	}
	for _, gid := range groupIDs {
		values.Add("group_id", strconv.FormatInt(int64(gid), 10))
	}
	if page > 1 {
		values.Set("page", strconv.FormatInt(page, 10))
	}
	encoded := values.Encode()
	if encoded == "" {
		return ""
	}
	return "?" + encoded
}

func productsEventsQS(onlyMissing bool, nameQuery string, groupIDs []products.GroupID, page int64) string {
	listQS := productsListQS(onlyMissing, nameQuery, groupIDs, page)
	eventsQS := "/events?topic=products-list"
	if listQS != "" {
		eventsQS = "/events?topic=products-list&" + listQS[1:]
	}
	return eventsQS
}

func formatQty(v products.Quantity) string {
	return v.String()
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func groupSelected(selected []products.GroupID, id products.GroupID) bool {
	for _, gid := range selected {
		if gid == id {
			return true
		}
	}
	return false
}

func productsTitle(onlyMissing bool) string {
	if onlyMissing {
		return "Braki / niski stan"
	}
	return "Wszystkie produkty"
}

func productsTabClass(active bool) string {
	if active {
		return "button"
	}
	return "button secondary"
}

func productRowClass(p products.Product) string {
	if p.IsMissing() || p.Quantity <= p.MinQuantity {
		return "warn"
	}
	return ""
}

func productIconKey(p products.Product) string {
	iconKey := strings.TrimSpace(string(p.IconKey))
	if iconKey == "" {
		return "cart"
	}
	return iconKey
}

func productCurrentGroupLabel(p products.Product) string {
	groupName := strings.TrimSpace(p.GroupName)
	if groupName == "" {
		return "(brak grupy)"
	}
	return groupName
}

func productGroupPostURL(p products.Product, listQS string) string {
	return "/products/" + strconv.FormatInt(int64(p.ID), 10) + "/group" + listQS
}

func shoppingItemID(item shoppinglist.Item) string {
	return strconv.FormatInt(int64(item.ID), 10)
}

func shoppingNextDone(done bool) string {
	if done {
		return "0"
	}
	return "1"
}

func shoppingNameClass(done bool) string {
	if done {
		return "sl-name done"
	}
	return "sl-name"
}

func shoppingDoneBtnClass(done bool) string {
	if done {
		return "icon-btn sl-done sl-done-on"
	}
	return "icon-btn sl-done sl-done-off"
}

// slItemDoneClass returns additional CSS class for done items.
func slItemDoneClass(done bool) string {
	if done {
		return " sl-item-done"
	}
	return ""
}

func shoppingIconSrc(item shoppinglist.Item) string {
	iconKey := strings.TrimSpace(item.IconKey)
	if iconKey != "" {
		return "/static/icons/" + iconKey + ".svg"
	}
	return "/icons/auto?name=" + url.QueryEscape(item.Name)
}

func shoppingGroupLabel(item shoppinglist.Item) string {
	return strings.TrimSpace(item.GroupName)
}

func normalizedUnit(u products.Unit) string {
	v := strings.TrimSpace(string(u))
	if v == "" {
		return string(products.UnitPiece)
	}
	return v
}

// quantityStep returns the HTML step attribute value for quantity inputs.
// Returns "1" for integer-only products, "0.1" otherwise.
func quantityStep(integerOnly bool) string {
	if integerOnly {
		return "1"
	}
	return "0.1"
}

// shoppingItemStep returns the HTML step attribute for shopping list item quantity input.
func shoppingItemStep(item shoppinglist.Item) string {
	return quantityStep(item.IntegerOnly)
}

// shoppingItemMin returns the HTML min attribute for shopping list item quantity input.
func shoppingItemMin(item shoppinglist.Item) string {
	return quantityMin(item.IntegerOnly)
}

// quantityMin returns the HTML min attribute value for quantity inputs.
// Returns "1" for integer-only products, "0.1" otherwise.
func quantityMin(integerOnly bool) string {
	if integerOnly {
		return "1"
	}
	return "0.1"
}

// groupProducts groups products by their group name.
// Products without a group are placed in a group with empty name.
// The order of products within each group is preserved from input.
func groupProducts(prods []products.Product) []ProductGroup {
	if len(prods) == 0 {
		return nil
	}

	// Use a slice to maintain order, map for lookup
	var groups []ProductGroup
	groupIndex := make(map[string]int)

	for _, p := range prods {
		groupName := p.GroupName
		if idx, exists := groupIndex[groupName]; exists {
			groups[idx].Products = append(groups[idx].Products, p)
		} else {
			groupIndex[groupName] = len(groups)
			groups = append(groups, ProductGroup{
				Name:     groupName,
				Products: []products.Product{p},
			})
		}
	}

	return groups
}

// groupShoppingItems groups shopping list items by their group name.
// Items without a group are placed in a group with empty name.
// The order of items within each group is preserved from input.
func groupShoppingItems(items []shoppinglist.Item) []ShoppingItemGroup {
	if len(items) == 0 {
		return nil
	}

	// Use a slice to maintain order, map for lookup
	var groups []ShoppingItemGroup
	groupIndex := make(map[string]int)

	for _, item := range items {
		groupName := item.GroupName
		if idx, exists := groupIndex[groupName]; exists {
			groups[idx].Items = append(groups[idx].Items, item)
		} else {
			groupIndex[groupName] = len(groups)
			groups = append(groups, ShoppingItemGroup{
				Name:  groupName,
				Items: []shoppinglist.Item{item},
			})
		}
	}

	return groups
}

// groupNameDisplay returns a display-friendly group name.
// Returns "(brak grupy)" for empty group names.
func groupNameDisplay(name string) string {
	if name == "" {
		return "(brak grupy)"
	}
	return name
}
