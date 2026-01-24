package migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	appmigrations "shopping/migrations"
)

func Up(db *sql.DB) error {
	src, err := iofs.New(appmigrations.FS, ".")
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "sqlite3", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		version, dirty, verr := m.Version()
		if verr != nil && !errors.Is(verr, migrate.ErrNilVersion) {
			return fmt.Errorf("migration failed: %w (could not get version: %v)", err, verr)
		}
		return fmt.Errorf("migration failed at version %d (dirty=%v): %w", version, dirty, err)
	}
	if err := ensureShoppingListDoneAt(db); err != nil {
		return err
	}
	return ensureUnitForms(db)
}

func ensureShoppingListDoneAt(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "PRAGMA table_info(shopping_list_items)")
	if err != nil {
		return err
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
			return err
		}
		if name == "done_at" {
			hasDoneAt = true
			break
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if hasDoneAt {
		return nil
	}

	if _, err := db.ExecContext(ctx, "ALTER TABLE shopping_list_items ADD COLUMN done_at TIMESTAMP NULL"); err != nil {
		return fmt.Errorf("add shopping_list_items.done_at: %w", err)
	}
	if _, err := db.ExecContext(ctx, "UPDATE shopping_list_items SET done_at = created_at WHERE done = 1 AND done_at IS NULL"); err != nil {
		return fmt.Errorf("backfill shopping_list_items.done_at: %w", err)
	}
	if _, err := db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_shopping_list_done_at ON shopping_list_items(done, done_at)"); err != nil {
		return fmt.Errorf("create idx_shopping_list_done_at: %w", err)
	}
	return nil
}

func ensureUnitForms(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "PRAGMA table_info(units)")
	if err != nil {
		return err
	}
	defer rows.Close()

	var hasSingular bool
	var hasPlural bool
	for rows.Next() {
		var cid int
		var name string
		var colType string
		var notNull int
		var dflt sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dflt, &pk); err != nil {
			return err
		}
		switch name {
		case "singular":
			hasSingular = true
		case "plural":
			hasPlural = true
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if !hasSingular {
		if _, err := db.ExecContext(ctx, "ALTER TABLE units ADD COLUMN singular TEXT NOT NULL DEFAULT ''"); err != nil {
			return fmt.Errorf("add units.singular: %w", err)
		}
	}
	if !hasPlural {
		if _, err := db.ExecContext(ctx, "ALTER TABLE units ADD COLUMN plural TEXT NOT NULL DEFAULT ''"); err != nil {
			return fmt.Errorf("add units.plural: %w", err)
		}
	}

	if _, err := db.ExecContext(ctx, "UPDATE products SET quantity_unit = 'sztuk' WHERE lower(quantity_unit) = 'sztuka'"); err != nil {
		return fmt.Errorf("normalize products.unit: %w", err)
	}
	if _, err := db.ExecContext(ctx, "UPDATE shopping_list_items SET quantity_unit = 'sztuk' WHERE lower(quantity_unit) = 'sztuka'"); err != nil {
		return fmt.Errorf("normalize shopping_list_items.unit: %w", err)
	}
	if _, err := db.ExecContext(ctx, "DELETE FROM units WHERE lower(name) = 'sztuka' AND EXISTS (SELECT 1 FROM units WHERE name = 'sztuk')"); err != nil {
		return fmt.Errorf("cleanup units.sztuka: %w", err)
	}
	if _, err := db.ExecContext(ctx, "UPDATE units SET name = 'sztuk' WHERE lower(name) = 'sztuka'"); err != nil {
		return fmt.Errorf("rename units.sztuka: %w", err)
	}
	if _, err := db.ExecContext(ctx, "INSERT OR IGNORE INTO units (name, display_order, singular, plural) VALUES ('sztuk', 1, 'sztuka', 'sztuk')"); err != nil {
		return fmt.Errorf("insert units.sztuk: %w", err)
	}

	if _, err := db.ExecContext(ctx, "UPDATE units SET singular = name, plural = name WHERE singular = '' AND plural = ''"); err != nil {
		return fmt.Errorf("default units forms: %w", err)
	}
	if _, err := db.ExecContext(ctx, "UPDATE units SET singular = 'sztuka', plural = 'sztuk' WHERE name = 'sztuk'"); err != nil {
		return fmt.Errorf("update units.sztuk forms: %w", err)
	}
	if _, err := db.ExecContext(ctx, "UPDATE units SET singular = 'opakowanie', plural = 'opakowania' WHERE name = 'opakowanie'"); err != nil {
		return fmt.Errorf("update units.opakowanie forms: %w", err)
	}
	if _, err := db.ExecContext(ctx, "UPDATE units SET singular = 'pęczek', plural = 'pęczki' WHERE name = 'pęczek'"); err != nil {
		return fmt.Errorf("update units.peczek forms: %w", err)
	}

	return nil
}
