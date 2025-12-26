package migrator

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestUpEnsuresShoppingListDoneAt(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite", "file:memdb1?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	// Simulate an inconsistent DB: schema_migrations claims we're on v4, but the table is missing `done_at`.
	_, err = db.Exec(`
		CREATE TABLE schema_migrations (version uint64, dirty bool);
		CREATE UNIQUE INDEX version_unique ON schema_migrations (version);
		INSERT INTO schema_migrations(version, dirty) VALUES (1, 0);
		CREATE TABLE shopping_list_items (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  product_id INTEGER NULL,
		  name TEXT NOT NULL,
		  quantity_value REAL NOT NULL DEFAULT 1,
		  quantity_unit TEXT NOT NULL DEFAULT 'sztuk',
		  done INTEGER NOT NULL DEFAULT 0,
		  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("seed schema: %v", err)
	}

	if err := Up(db); err != nil {
		t.Fatalf("Up: %v", err)
	}

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "PRAGMA table_info(shopping_list_items)")
	if err != nil {
		t.Fatalf("pragma: %v", err)
	}
	defer rows.Close()

	hasDoneAt := false
	for rows.Next() {
		var cid int
		var name string
		var colType string
		var notNull int
		var dflt sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dflt, &pk); err != nil {
			t.Fatalf("scan: %v", err)
		}
		if name == "done_at" {
			hasDoneAt = true
			break
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows: %v", err)
	}
	if !hasDoneAt {
		t.Fatalf("expected shopping_list_items.done_at to exist")
	}
}
