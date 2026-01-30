# Domain Concepts

This document describes the core domain concepts: products, groups, shopping list items, and their relationships.

---

## Products

Tracked supplies in the home inventory system.

### Attributes

| Field | Type | Description |
|-------|------|-------------|
| `id` | Integer | Unique identifier |
| `name` | String | Product name (Polish, unique) |
| `group_id` | Integer | Foreign key to groups table |
| `icon_key` | String | Icon identifier (e.g., "carrot", "milk") |
| `quantity_value` | Float | Current quantity (0 = missing/out of stock) |
| `quantity_unit` | String | Unit: `sztuk`, `kg`, `litr`, `gramy`, `opakowanie` |
| `min_quantity_value` | Float | Minimum threshold for "low stock" warning |
| `integer_only` | Boolean | 1 if item can't have decimals (eggs, packages) |
| `created_at` | Timestamp | Creation time |
| `updated_at` | Timestamp | Last update time |

### Missing Status

Products with `quantity_value = 0` are considered **missing** (out of stock). There is no separate `missing` boolean column - the status is derived from quantity.

```go
func (p Product) IsMissing() bool {
    return p.QuantityValue == 0
}
```

### Units (Polish)

| Unit | Usage |
|------|-------|
| `sztuk` | Pieces (default) - eggs, apples |
| `kg` | Kilograms - vegetables, meat |
| `litr` | Liters - milk, juice |
| `gramy` | Grams - cheese, small quantities |
| `opakowanie` | Package - boxed items |

### Auto-Icon Resolution

Products get icons automatically based on name pattern matching:

1. Check `product_icon_rules` table for matching substring
2. Higher priority rules checked first
3. First match wins
4. Default to "cart" if no match

**Example rules**:
```sql
INSERT INTO product_icon_rules (match_substring, icon_key, priority) VALUES
('marchew', 'carrot', 100),
('mleko', 'milk', 100);
```

When adding "marchewka", it matches "marchew" rule → gets carrot icon.

### Schema Reference

[internal/domain/products/entities.go](../internal/domain/products/entities.go)
[migrations/0001_init.up.sql](../migrations/0001_init.up.sql) (products table)

---

## Groups

Product categories/organizational groups.

### Attributes

| Field | Type | Description |
|-------|------|-------------|
| `id` | Integer | Unique identifier |
| `name` | String | Group name (Polish, unique) |

### Common Groups

| Group Name | Translation |
|------------|-------------|
| `warzywa` | Vegetables |
| `owoce` | Fruits |
| `nabiał` | Dairy |
| `mięso` | Meat |
| `pieczywo` | Bread/Bakery |
| `napoje` | Beverages |
| `chemia` | Cleaning supplies |

### Views

The `v_groups` view provides group information joined with product counts:
```sql
CREATE VIEW v_groups AS
SELECT g.*, COUNT(p.id) as product_count
FROM groups g
LEFT JOIN products p ON p.group_id = g.id
GROUP BY g.id;
```

### Schema Reference

[migrations/0001_init.up.sql](../migrations/0001_init.up.sql) (groups table)

---

## Shopping List

Items to purchase, optionally linked to existing products.

### Attributes

| Field | Type | Description |
|-------|------|-------------|
| `id` | Integer | Unique identifier |
| `name` | String | Item name (Polish) |
| `product_id` | Integer? | Optional FK to products table |
| `quantity_value` | Float | Quantity to buy |
| `quantity_unit` | String | Unit (same options as products) |
| `done` | Boolean | Marked as purchased |
| `done_at` | Timestamp? | When marked done |
| `created_at` | Timestamp | Creation time |

### Auto-Linking to Products

When adding an item to the shopping list:
1. Search for existing product with matching name (case-insensitive)
2. If found, link via `product_id`
3. Copy product's icon and unit

This enables auto-inventory updates when marking done.

### Auto-Inventory Update

When marking a shopping item as **done**:
1. If linked to product → increment product quantity
2. Shopping item stays in list with `done=true`
3. Auto-cleanup removes done items after 6 hours

### Cleanup

A background job runs periodically to remove shopping items where:
- `done = true`
- `done_at < now() - 6 hours`

### Schema Reference

[internal/domain/shoppinglist/types.go](../internal/domain/shoppinglist/types.go)
[migrations/0001_init.up.sql](../migrations/0001_init.up.sql) (shopping_list table)

---

## Relationships

```
┌─────────────┐       ┌─────────────┐       ┌─────────────────┐
│   Groups    │       │   Products  │       │  Shopping List  │
├─────────────┤       ├─────────────┤       ├─────────────────┤
│ id          │◄──────┤ group_id    │       │ id              │
│ name        │       │ name        │◄──────┤ product_id (opt)│
└─────────────┘       │ quantity_   │       │ name            │
                      │   value     │       │ quantity_value  │
                      │ icon_key    │◄──────┤ (copied)        │
                      └─────────────┘       │ done            │
                                            │ done_at         │
                                            └─────────────────┘
```

### Business Rules

1. **Product names are unique** - Can't have two "Milk" products
2. **Group names are unique** - Can't have duplicate group names
3. **Zero quantity = missing** - Products with no stock are flagged
4. **Shopping items link by name** - Auto-matching enables inventory updates
5. **Done items auto-expire** - Cleaned up after 6 hours
6. **Icon rules use substring matching** - "marchewka" matches "marchew" rule

---

## Domain Services

### Product Service

**File**: [internal/domain/products/service.go](../internal/domain/products/service.go)

**Responsibilities**:
- Create, update, delete products
- Validate product data
- Handle icon resolution
- Maintain product groups

**Key methods**:
```go
Create(ctx context.Context, cmd CreateProductCommand) error
Update(ctx context.Context, id int64, cmd UpdateProductCommand) error
Delete(ctx context.Context, id int64) error
UpdateQuantity(ctx context.Context, id int64, delta float64) error
```

### Shopping List Service

**File**: [internal/domain/shoppinglist/service.go](../internal/domain/shoppinglist/service.go)

**Responsibilities**:
- Add items to shopping list
- Mark items as done
- Auto-link to products
- Auto-update inventory
- Cleanup old done items

**Key methods**:
```go
AddItem(ctx context.Context, cmd AddItemCommand) error
MarkDone(ctx context.Context, id int64) error
MarkUndone(ctx context.Context, id int64) error
DeleteItem(ctx context.Context, id int64) error
CleanupDone(ctx context.Context) error
```

### Admin Service

**File**: [internal/domain/admin/ports.go](../internal/domain/admin/ports.go)

**Responsibilities**:
- Database optimization
- Maintenance operations

---

## Validation Rules

### Product Validation

**File**: [internal/domain/products/validation.go](../internal/domain/products/validation.go)

| Rule | Constraint |
|------|-----------|
| Name | Required, 1-100 characters |
| Name | Unique (database constraint) |
| Quantity | Non-negative |
| Min Quantity | Non-negative, <= quantity |
| Unit | Must be valid unit type |

### Shopping Item Validation

**File**: [internal/domain/shoppinglist/validation.go](../internal/domain/shoppinglist/validation.go)

| Rule | Constraint |
|------|-----------|
| Name | Required, 1-100 characters |
| Quantity | Positive number |
| Unit | Must be valid unit type |
