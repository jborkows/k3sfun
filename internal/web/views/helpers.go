package views

import (
	"io"
	"net/url"
	"shopping/internal/domain/products"
	"shopping/internal/domain/shoppinglist"
	"sort"
	"strconv"
	"strings"
)

func boolToQuery(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func productsListQS(onlyMissing bool, nameQuery string, groups []products.Group, groupIDs []products.GroupID, page int64) string {
	values := url.Values{}
	if onlyMissing {
		values.Set("missing", "1")
	}
	nameQuery = strings.TrimSpace(nameQuery)
	if nameQuery != "" {
		values.Set("q", nameQuery)
	}
	if names := products.GroupIDsToNames(groups, groupIDs); len(names) > 0 {
		values.Set("groups", strings.Join(names, ","))
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

func productsEventsQS(onlyMissing bool, nameQuery string, groups []products.Group, groupIDs []products.GroupID, page int64) string {
	listQS := productsListQS(onlyMissing, nameQuery, groups, groupIDs, page)
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

// groupNameByID returns the name of a group given its ID.
// Returns empty string if the group is not found.
func groupNameByID(groups []products.Group, id products.GroupID) string {
	for _, g := range groups {
		if g.ID == id {
			return g.Name
		}
	}
	return ""
}

// productsListQSWithoutGroup returns a query string without a specific group ID.
// Used for filter chip removal.
func productsListQSWithoutGroup(onlyMissing bool, nameQuery string, groups []products.Group, groupIDs []products.GroupID, page int64, excludeGroupID products.GroupID) string {
	var filtered []products.GroupID
	for _, gid := range groupIDs {
		if gid != excludeGroupID {
			filtered = append(filtered, gid)
		}
	}
	return productsListQS(onlyMissing, nameQuery, groups, filtered, page)
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

func brakBtnClass(p products.Product) string {
	if p.IsMissing() || p.Quantity <= p.MinQuantity {
		return " brak-btn-active"
	}
	return ""
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

	ordered := make([]shoppinglist.Item, len(items))
	copy(ordered, items)
	sortShoppingItems(ordered)

	// Use a slice to maintain order, map for lookup
	var groups []ShoppingItemGroup
	groupIndex := make(map[string]int)

	for _, item := range ordered {
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

func sortShoppingItems(items []shoppinglist.Item) {
	sort.SliceStable(items, func(i, j int) bool {
		a := items[i]
		b := items[j]
		if a.GroupOrder != b.GroupOrder {
			return a.GroupOrder < b.GroupOrder
		}
		if groupCmp := strings.Compare(a.GroupName, b.GroupName); groupCmp != 0 {
			return groupCmp < 0
		}
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) < 0
	})
}

func RenderShoppingListExport(w io.Writer, items []shoppinglist.Item) error {
	for _, item := range items {
		line := strings.TrimSpace(item.Name)
		if line == "" {
			continue
		}
		if _, err := io.WriteString(w, line+"\t"+formatQty(item.Quantity)+"\t"+normalizedUnit(item.Unit)+"\n"); err != nil {
			return err
		}
	}
	return nil
}

// groupNameDisplay returns a display-friendly group name.
// Returns "(brak grupy)" for empty group names.
func groupNameDisplay(name string) string {
	if name == "" {
		return "(brak grupy)"
	}
	return name
}
