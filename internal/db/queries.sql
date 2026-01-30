-- name: ListGroups :many
SELECT id, name
FROM v_groups
ORDER BY display_order, name;

-- name: CreateGroup :one
INSERT INTO groups(name) VALUES (?)
RETURNING id;

-- name: ListProductsAll :many
SELECT
  p.id,
  p.name,
  p.icon_key,
  p.group_id,
  p.group_name,
  p.quantity_value,
  p.quantity_unit,
  p.updated_at
FROM v_products p
ORDER BY p.name;

-- name: ListProductsMissingOrLow :many
SELECT
  p.id,
  p.name,
  p.icon_key,
  p.group_id,
  p.group_name,
  p.quantity_value,
  p.quantity_unit,
  p.updated_at
FROM v_products p
WHERE p.quantity_value = 0
ORDER BY p.name;

-- name: ListProductsFiltered :many
SELECT
  p.id,
  p.name,
  p.icon_key,
  p.group_id,
  p.group_name,
  p.quantity_value,
  p.quantity_unit,
  p.updated_at
FROM v_products p
WHERE
  (? = 0 OR p.quantity_value = 0)
  AND (? = '' OR lower(p.name) LIKE '%' || lower(?) || '%')
  AND (? = 0 OR p.group_id IN (sqlc.slice('group_ids')))
ORDER BY COALESCE(lower(p.group_name), 'zzz'), lower(p.name)
LIMIT ?
OFFSET ?;

-- name: CountProductsFiltered :one
SELECT COUNT(*)
FROM v_products p
WHERE
  (? = 0 OR p.quantity_value = 0)
  AND (? = '' OR lower(p.name) LIKE '%' || lower(?) || '%')
  AND (? = 0 OR p.group_id IN (sqlc.slice('group_ids')));

-- name: CreateProduct :one
INSERT INTO products(name, icon_key, group_id, quantity_value, quantity_unit, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id;

-- name: SetProductQuantity :exec
UPDATE products
SET quantity_value = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: AddProductQuantity :exec
UPDATE products
SET quantity_value = quantity_value + ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SetProductUnit :exec
UPDATE products
SET quantity_unit = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;



-- name: SetProductGroup :exec
UPDATE products
SET group_id = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: ListShoppingListItems :many
SELECT
  sli.id,
  sli.product_id,
  sli.name,
  COALESCE(p.icon_key, '') AS icon_key,
  COALESCE(g.name, '') AS group_name,
  COALESCE(g.display_order, 999) AS group_order,
  sli.quantity_value,
  CASE WHEN lower(sli.quantity_unit) = 'sztuka' THEN 'sztuk' ELSE sli.quantity_unit END AS quantity_unit,
  COALESCE(u.singular, CASE WHEN lower(sli.quantity_unit) = 'sztuka' THEN 'sztuk' ELSE sli.quantity_unit END) AS unit_singular,
  COALESCE(u.plural, CASE WHEN lower(sli.quantity_unit) = 'sztuka' THEN 'sztuk' ELSE sli.quantity_unit END) AS unit_plural,
  sli.done,
  sli.created_at
FROM shopping_list_items sli
LEFT JOIN products p ON p.id = sli.product_id
LEFT JOIN groups g ON g.id = p.group_id
LEFT JOIN units u ON u.name = CASE WHEN lower(sli.quantity_unit) = 'sztuka' THEN 'sztuk' ELSE sli.quantity_unit END
ORDER BY sli.done ASC, COALESCE(g.display_order, 999), COALESCE(lower(g.name), 'zzz'), lower(sli.name);

-- name: GetShoppingListItem :one
SELECT
  sli.id,
  sli.product_id,
  sli.name,
  COALESCE(p.icon_key, '') AS icon_key,
  COALESCE(g.name, '') AS group_name,
  COALESCE(g.display_order, 999) AS group_order,
  sli.quantity_value,
  CASE WHEN lower(sli.quantity_unit) = 'sztuka' THEN 'sztuk' ELSE sli.quantity_unit END AS quantity_unit,
  COALESCE(u.singular, CASE WHEN lower(sli.quantity_unit) = 'sztuka' THEN 'sztuk' ELSE sli.quantity_unit END) AS unit_singular,
  COALESCE(u.plural, CASE WHEN lower(sli.quantity_unit) = 'sztuka' THEN 'sztuk' ELSE sli.quantity_unit END) AS unit_plural,
  sli.done,
  sli.created_at
FROM shopping_list_items sli
LEFT JOIN products p ON p.id = sli.product_id
LEFT JOIN groups g ON g.id = p.group_id
LEFT JOIN units u ON u.name = CASE WHEN lower(sli.quantity_unit) = 'sztuka' THEN 'sztuk' ELSE sli.quantity_unit END
WHERE sli.id = ?;

-- name: AddShoppingListItemByName :exec
INSERT OR IGNORE INTO shopping_list_items(product_id, name, quantity_value, quantity_unit, done)
VALUES (NULL, ?, ?, ?, 0);

-- name: AddShoppingListItemByProductID :exec
INSERT OR IGNORE INTO shopping_list_items(product_id, name, quantity_value, quantity_unit, done)
SELECT p.id, p.name, 1, p.quantity_unit, 0
FROM products p
WHERE p.id = ?;

-- name: SetShoppingListItemDone :exec
UPDATE shopping_list_items
SET
  done = ?,
  done_at = CASE WHEN ? = 1 THEN CURRENT_TIMESTAMP ELSE NULL END
WHERE id = ?;

-- name: CleanupShoppingListDoneBefore :exec
DELETE FROM shopping_list_items
WHERE
  done = 1
  AND COALESCE(done_at, created_at) < ?;

-- name: SetShoppingListItemQuantity :exec
UPDATE shopping_list_items
SET quantity_value = ?, quantity_unit = ?
WHERE id = ?;

-- name: DeleteShoppingListItem :exec
DELETE FROM shopping_list_items
WHERE id = ?;

-- name: LinkShoppingListItemToProduct :exec
UPDATE shopping_list_items
SET product_id = ?, name = ?
WHERE id = ?;

-- name: FindProductIDByName :one
SELECT id
FROM products
WHERE lower(name) = lower(?)
LIMIT 1;

-- name: SuggestProductsByName :many
SELECT
  id,
  name,
  icon_key,
  quantity_unit
FROM products
WHERE lower(name) LIKE '%' || lower(?) || '%'
ORDER BY lower(name)
LIMIT ?;

-- name: ResolveProductIconKeyByName :one
SELECT icon_key
FROM product_icon_rules
WHERE lower(?) LIKE '%' || lower(match_substring) || '%'
ORDER BY priority DESC, id DESC
LIMIT 1;

-- name: ListUnits :many
SELECT name
FROM units
ORDER BY display_order, name;
