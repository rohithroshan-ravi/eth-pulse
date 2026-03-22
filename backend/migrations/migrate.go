package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rohithroshan-ravi/eth-pulse/backend/config"
)

//go:generate ls migrations/*.sql > /dev/null 2>&1 || true

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run migrations/migrate.go [up|down|status]")
		os.Exit(1)
	}

	command := os.Args[1]

	// Connect to database
	config.ConnectDB()
	defer config.DB.Close()

	ctx := context.Background()

	// Create migration history table if it doesn't exist
	if err := createMigrationHistoryTable(ctx); err != nil {
		log.Fatalf("Failed to create migration history table: %v", err)
	}

	switch command {
	case "up":
		if err := migrateUp(ctx); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("✅ Migrations applied successfully")

	case "down":
		fmt.Println("⚠️ Migration down not yet implemented")

	case "status":
		if err := showMigrationStatus(ctx); err != nil {
			log.Fatalf("Failed to show migration status: %v", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func createMigrationHistoryTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS public.migration_history (
		id SERIAL PRIMARY KEY,
		migration_name text NOT NULL UNIQUE,
		executed_at timestamptz NOT NULL DEFAULT now()
	);
	`
	_, err := config.DB.Exec(ctx, query)
	return err
}

func migrateUp(ctx context.Context) error {
	// Get list of migration files
	migrations, err := getMigrationFiles()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		// Check if migration already applied
		applied, err := isMigrationApplied(ctx, migration)
		if err != nil {
			return err
		}

		if applied {
			fmt.Printf("⏭️  Skipping %s (already applied)\n", migration)
			continue
		}

		// Read migration file
		content, err := os.ReadFile(filepath.Join("migrations", migration))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", migration, err)
		}

		// Execute migration
		fmt.Printf("⬆️  Applying %s...\n", migration)
		if err := executeMigration(ctx, string(content), migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration, err)
		}

		fmt.Printf("✅ Applied %s\n", migration)
	}

	return nil
}

func getMigrationFiles() ([]string, error) {
	files, err := os.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, file.Name())
		}
	}

	sort.Strings(migrations)
	return migrations, nil
}

func isMigrationApplied(ctx context.Context, name string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM public.migration_history WHERE migration_name = $1)"
	err := config.DB.QueryRow(ctx, query, name).Scan(&exists)
	return exists, err
}

func executeMigration(ctx context.Context, sql, name string) error {
	// Begin transaction
	tx, err := config.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Execute migration
	if _, err := tx.Exec(ctx, sql); err != nil {
		return err
	}

	// Record migration
	query := "INSERT INTO public.migration_history (migration_name) VALUES ($1)"
	if _, err := tx.Exec(ctx, query, name); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func showMigrationStatus(ctx context.Context) error {
	// Get all migration files
	files, err := getMigrationFiles()
	if err != nil {
		return err
	}

	fmt.Println("\n📋 Migration Status:")
	fmt.Println("-------------------")

	// Get applied migrations
	rows, err := config.DB.Query(ctx, "SELECT migration_name, executed_at FROM public.migration_history ORDER BY executed_at")
	if err != nil {
		return err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		var executedAt interface{}
		if err := rows.Scan(&name, &executedAt); err != nil {
			return err
		}
		applied[name] = true
		fmt.Printf("✅ %s - Applied at %v\n", name, executedAt)
	}

	// Show pending migrations
	pending := false
	for _, file := range files {
		if !applied[file] {
			if !pending {
				fmt.Println("\n⏳ Pending migrations:")
				pending = true
			}
			fmt.Printf("⏭️  %s\n", file)
		}
	}

	if len(applied) == 0 {
		fmt.Println("No migrations applied yet")
	}

	return nil
}
