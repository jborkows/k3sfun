package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"shopping/internal/domain/products"
	"shopping/internal/migrator"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Ping: %v", err)
	}
	return db
}

// setupCleanDB runs migrations and truncates all seeded data to start with a clean slate.
func setupCleanDB(t *testing.T, db *sql.DB) {
	t.Helper()

	if err := migrator.Up(db); err != nil {
		t.Fatalf("migrator.Up: %v", err)
	}

	// Truncate seeded data to start fresh
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tables := []string{"shopping_list_items", "products", "groups", "product_icon_rules"}
	for _, table := range tables {
		if _, err := db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			t.Fatalf("truncate %s: %v", table, err)
		}
	}
}

func TestRepo_ListProducts_FilteringAndPaging(t *testing.T) {
	db := openTestDB(t)
	setupCleanDB(t, db)

	r := NewRepo(db)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create test groups
	vegetablesID, err := r.CreateGroup(ctx, "warzywa")
	if err != nil {
		t.Fatalf("CreateGroup(warzywa): %v", err)
	}
	fruitsID, err := r.CreateGroup(ctx, "owoce")
	if err != nil {
		t.Fatalf("CreateGroup(owoce): %v", err)
	}
	dairyID, err := r.CreateGroup(ctx, "nabiał")
	if err != nil {
		t.Fatalf("CreateGroup(nabiał): %v", err)
	}

	// Create test products
	testProducts := []products.NewProduct{
		{Name: "marchewka", IconKey: "carrot", GroupID: &vegetablesID, Quantity: 1.5, Unit: products.UnitKG},
		{Name: "ziemniaki", IconKey: "potato", GroupID: &vegetablesID, Quantity: 2.0, Unit: products.UnitKG},
		{Name: "jabłka", IconKey: "apple", GroupID: &fruitsID, Quantity: 0, Unit: products.UnitPiece},
		{Name: "mleko", IconKey: "milk", GroupID: &dairyID, Quantity: 1.0, Unit: products.UnitLiter},
		{Name: "śmietana", IconKey: "sour-cream", GroupID: &dairyID, Quantity: 0, Unit: products.UnitLiter},
	}

	createdIDs := make([]products.ProductID, 0, len(testProducts))
	for _, p := range testProducts {
		id, err := r.CreateProduct(ctx, p)
		if err != nil {
			t.Fatalf("CreateProduct(%s): %v", p.Name, err)
		}
		createdIDs = append(createdIDs, id)
	}

	t.Run("empty filter returns all products", func(t *testing.T) {
		all, err := r.ListProducts(ctx, products.ProductFilter{
			Limit: products.MaxProductsPageSize,
		})
		if err != nil {
			t.Fatalf("ListProducts: %v", err)
		}
		if len(all) != len(testProducts) {
			t.Errorf("expected %d products, got %d", len(testProducts), len(all))
		}
	})

	t.Run("filter by single group", func(t *testing.T) {
		list, err := r.ListProducts(ctx, products.ProductFilter{
			GroupIDs: []products.GroupID{vegetablesID},
			Limit:    products.MaxProductsPageSize,
		})
		if err != nil {
			t.Fatalf("ListProducts: %v", err)
		}
		if len(list) != 2 {
			t.Errorf("expected 2 vegetables, got %d", len(list))
		}
		for _, p := range list {
			if p.GroupID == nil || *p.GroupID != vegetablesID {
				t.Errorf("product %q should be in vegetables group", p.Name)
			}
		}
	})

	t.Run("filter by multiple groups", func(t *testing.T) {
		list, err := r.ListProducts(ctx, products.ProductFilter{
			GroupIDs: []products.GroupID{vegetablesID, fruitsID},
			Limit:    products.MaxProductsPageSize,
		})
		if err != nil {
			t.Fatalf("ListProducts: %v", err)
		}
		if len(list) != 3 {
			t.Errorf("expected 3 products (2 vegetables + 1 fruit), got %d", len(list))
		}
	})

	t.Run("filter by name query", func(t *testing.T) {
		list, err := r.ListProducts(ctx, products.ProductFilter{
			NameQuery: "marchew",
			Limit:     products.MaxProductsPageSize,
		})
		if err != nil {
			t.Fatalf("ListProducts: %v", err)
		}
		if len(list) != 1 {
			t.Errorf("expected 1 product matching 'marchew', got %d", len(list))
		}
		if len(list) > 0 && list[0].Name != "marchewka" {
			t.Errorf("expected 'marchewka', got %q", list[0].Name)
		}
	})

	t.Run("filter by name with Polish diacritics", func(t *testing.T) {
		// Test uppercase Polish letter in search
		list, err := r.ListProducts(ctx, products.ProductFilter{
			NameQuery: "Śmie",
			Limit:     products.MaxProductsPageSize,
		})
		if err != nil {
			t.Fatalf("ListProducts: %v", err)
		}
		if len(list) != 1 {
			t.Errorf("expected 1 product matching 'Śmie', got %d", len(list))
		}
		if len(list) > 0 && list[0].Name != "śmietana" {
			t.Errorf("expected 'śmietana', got %q", list[0].Name)
		}
	})

	t.Run("filter by name and group combined", func(t *testing.T) {
		list, err := r.ListProducts(ctx, products.ProductFilter{
			NameQuery: "ml",
			GroupIDs:  []products.GroupID{dairyID},
			Limit:     products.MaxProductsPageSize,
		})
		if err != nil {
			t.Fatalf("ListProducts: %v", err)
		}
		if len(list) != 1 {
			t.Errorf("expected 1 product, got %d", len(list))
		}
		if len(list) > 0 && list[0].Name != "mleko" {
			t.Errorf("expected 'mleko', got %q", list[0].Name)
		}
	})

	t.Run("filter missing or low quantity", func(t *testing.T) {
		list, err := r.ListProducts(ctx, products.ProductFilter{
			OnlyMissingOrLow: true,
			Limit:            products.MaxProductsPageSize,
		})
		if err != nil {
			t.Fatalf("ListProducts: %v", err)
		}
		// jabłka and śmietana have quantity=0, so they should be marked as missing
		if len(list) != 2 {
			t.Errorf("expected 2 missing/low products, got %d", len(list))
		}
	})

	t.Run("count products with filter", func(t *testing.T) {
		count, err := r.CountProducts(ctx, products.ProductFilter{
			GroupIDs: []products.GroupID{vegetablesID},
		})
		if err != nil {
			t.Fatalf("CountProducts: %v", err)
		}
		if count != 2 {
			t.Errorf("expected count=2, got %d", count)
		}
	})

	t.Run("pagination offset", func(t *testing.T) {
		// Get first 2 products
		first, err := r.ListProducts(ctx, products.ProductFilter{
			Limit:  2,
			Offset: 0,
		})
		if err != nil {
			t.Fatalf("ListProducts (first): %v", err)
		}
		// Get next 2 products
		second, err := r.ListProducts(ctx, products.ProductFilter{
			Limit:  2,
			Offset: 2,
		})
		if err != nil {
			t.Fatalf("ListProducts (second): %v", err)
		}

		if len(first) != 2 {
			t.Errorf("expected 2 products in first page, got %d", len(first))
		}
		if len(second) != 2 {
			t.Errorf("expected 2 products in second page, got %d", len(second))
		}

		// Ensure no overlap
		for _, p1 := range first {
			for _, p2 := range second {
				if p1.ID == p2.ID {
					t.Errorf("pagination overlap: product %q appears in both pages", p1.Name)
				}
			}
		}
	})

	t.Run("limit is clamped to MaxProductsPageSize", func(t *testing.T) {
		// This should not error even with very high limit
		_, err := r.ListProducts(ctx, products.ProductFilter{
			Limit:  9999,
			Offset: 0,
		})
		if err != nil {
			t.Fatalf("ListProducts (limit clamp): %v", err)
		}
	})
}

func TestRepo_SetProductQuantity_SyncsMissingFlag(t *testing.T) {
	db := openTestDB(t)
	setupCleanDB(t, db)

	r := NewRepo(db)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create a product with quantity > 0
	id, err := r.CreateProduct(ctx, products.NewProduct{
		Name:     "test product",
		IconKey:  "cart",
		Quantity: 5,
		Unit:     products.UnitPiece,
	})
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}

	// Set quantity to 0 - should set missing flag
	if err := r.SetProductQuantity(ctx, id, 0); err != nil {
		t.Fatalf("SetProductQuantity(0): %v", err)
	}

	// Verify product is now missing
	list, err := r.ListProducts(ctx, products.ProductFilter{OnlyMissingOrLow: true, Limit: 10})
	if err != nil {
		t.Fatalf("ListProducts: %v", err)
	}
	found := false
	for _, p := range list {
		if p.ID == id {
			found = true
			if !p.Missing {
				t.Errorf("product with quantity=0 should be missing")
			}
		}
	}
	if !found {
		t.Errorf("product should appear in missing/low list after setting quantity=0")
	}

	// Set quantity > 0 - should clear missing flag
	if err := r.SetProductQuantity(ctx, id, 5); err != nil {
		t.Fatalf("SetProductQuantity(5): %v", err)
	}

	list, err = r.ListProducts(ctx, products.ProductFilter{OnlyMissingOrLow: true, Limit: 10})
	if err != nil {
		t.Fatalf("ListProducts: %v", err)
	}
	for _, p := range list {
		if p.ID == id {
			t.Errorf("product with quantity>0 should not be in missing/low list")
		}
	}

	// Set quantity back to 0 - should set missing flag again
	if err := r.SetProductQuantity(ctx, id, 0); err != nil {
		t.Fatalf("SetProductQuantity(0): %v", err)
	}

	list, err = r.ListProducts(ctx, products.ProductFilter{OnlyMissingOrLow: true, Limit: 10})
	if err != nil {
		t.Fatalf("ListProducts: %v", err)
	}
	found = false
	for _, p := range list {
		if p.ID == id {
			found = true
			if !p.Missing {
				t.Errorf("product with quantity=0 should be missing")
			}
		}
	}
	if !found {
		t.Errorf("product should appear in missing/low list after setting quantity=0")
	}
}

func TestRepo_SuggestProductsByName_PolishDiacritics(t *testing.T) {
	db := openTestDB(t)
	setupCleanDB(t, db)

	r := NewRepo(db)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create products with Polish characters
	testProducts := []products.NewProduct{
		{Name: "śmietana", IconKey: "sour-cream", Quantity: 1, Unit: products.UnitLiter},
		{Name: "żurawina", IconKey: "cart", Quantity: 1, Unit: products.UnitKG},
		{Name: "ćwikła", IconKey: "beetroot", Quantity: 1, Unit: products.UnitPiece},
	}

	for _, p := range testProducts {
		if _, err := r.CreateProduct(ctx, p); err != nil {
			t.Fatalf("CreateProduct(%s): %v", p.Name, err)
		}
	}

	tests := []struct {
		query    string
		expected string
	}{
		{"Śmie", "śmietana"},
		{"śmie", "śmietana"},
		{"ŚMIE", "śmietana"},
		{"Żura", "żurawina"},
		{"żura", "żurawina"},
		{"Ćwi", "ćwikła"},
		{"ćwi", "ćwikła"},
	}

	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			suggestions, err := r.SuggestProductsByName(ctx, tc.query, 10)
			if err != nil {
				t.Fatalf("SuggestProductsByName(%q): %v", tc.query, err)
			}
			if len(suggestions) == 0 {
				t.Errorf("expected suggestion for %q, got none", tc.query)
				return
			}
			if suggestions[0].Name != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, suggestions[0].Name)
			}
		})
	}
}
