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
		return err
	}
	return ensureShoppingListDoneAt(db)
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
