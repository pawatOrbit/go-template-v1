package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/yourorg/go-api-template/config"
	"github.com/yourorg/go-api-template/core/pgdb"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  "Database migration commands to manage database schema changes",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Long:  "Apply all pending database migrations",
	RunE:  runMigrateUp,
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	Long:  "Rollback database migrations. Use --steps flag to specify number of migrations to rollback",
	RunE:  runMigrateDown,
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long:  "Show current migration status and pending migrations",
	RunE:  runMigrateStatus,
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration",
	Long:  "Create a new migration with up and down SQL files",
	Args:  cobra.ExactArgs(1),
	RunE:  runMigrateCreate,
}

var migrateForceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Force database to specific migration version",
	Long:  "Force set the database migration version (use with caution)",
	Args:  cobra.ExactArgs(1),
	RunE:  runMigrateForce,
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Long:  "Show the current database migration version",
	RunE:  runMigrateVersion,
}

var (
	migrateSteps int
	migrateAll   bool
)

func init() {
	// Add migrate command to root
	rootCmd.AddCommand(migrateCmd)

	// Add subcommands
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCmd.AddCommand(migrateForceCmd)
	migrateCmd.AddCommand(migrateVersionCmd)

	// Add flags
	migrateDownCmd.Flags().IntVar(&migrateSteps, "steps", 1, "Number of migrations to rollback")
	migrateDownCmd.Flags().BoolVar(&migrateAll, "all", false, "Rollback all migrations")
}

func getMigrationInstance() (*migrate.Migrate, error) {
	// Load configuration
	ctx := context.Background()
	if err := config.ResolveConfigFromFile(ctx, "config/config.local.yaml"); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	// Build database URL
	dbURL := buildDatabaseURL(cfg.Postgres.Write)

	// Get migrations path
	migrationsPath := "file://migrations"

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}

func buildDatabaseURL(dbConfig pgdb.PostgresConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)
}

func runMigrateUp(cmd *cobra.Command, args []string) error {
	m, err := getMigrationInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	fmt.Println("Running migrations...")

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("No migrations to run")
	} else {
		fmt.Println("Migrations completed successfully")
	}

	return nil
}

func runMigrateDown(cmd *cobra.Command, args []string) error {
	m, err := getMigrationInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	if migrateAll {
		fmt.Println("Rolling back all migrations...")
		err = m.Down()
	} else {
		fmt.Printf("Rolling back %d migration(s)...\n", migrateSteps)
		err = m.Steps(-migrateSteps)
	}

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("No migrations to rollback")
	} else {
		fmt.Println("Rollback completed successfully")
	}

	return nil
}

func runMigrateStatus(cmd *cobra.Command, args []string) error {
	m, err := getMigrationInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	if err == migrate.ErrNilVersion {
		fmt.Println("Database is not initialized")
		return nil
	}

	fmt.Printf("Current migration version: %d\n", version)
	if dirty {
		fmt.Println("Database is in dirty state - manual intervention may be required")
	} else {
		fmt.Println("Database is in clean state")
	}

	return nil
}

func runMigrateCreate(cmd *cobra.Command, args []string) error {
	name := args[0]
	timestamp := time.Now().Format("20060102150405")

	// Create migrations directory if it doesn't exist
	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(migrationsDir, 0755); err != nil {
			return fmt.Errorf("failed to create migrations directory: %w", err)
		}
	}

	// Create up migration file
	upFilename := fmt.Sprintf("%s_%s.up.sql", timestamp, name)
	upPath := filepath.Join(migrationsDir, upFilename)

	upContent := fmt.Sprintf("-- Migration: %s\n-- Created: %s\n-- Description: %s\n\n-- Add your up migration here\n",
		name, time.Now().Format(time.RFC3339), name)

	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	// Create down migration file
	downFilename := fmt.Sprintf("%s_%s.down.sql", timestamp, name)
	downPath := filepath.Join(migrationsDir, downFilename)

	downContent := fmt.Sprintf("-- Migration: %s (rollback)\n-- Created: %s\n-- Description: Rollback %s\n\n-- Add your down migration here\n",
		name, time.Now().Format(time.RFC3339), name)

	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	fmt.Printf("Created migration files:\n")
	fmt.Printf("  %s\n", upPath)
	fmt.Printf("  %s\n", downPath)

	return nil
}

func runMigrateForce(cmd *cobra.Command, args []string) error {
	versionStr := args[0]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return fmt.Errorf("invalid version number: %s", versionStr)
	}

	m, err := getMigrationInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	fmt.Printf("Forcing database to version %d...\n", version)

	err = m.Force(version)
	if err != nil {
		return fmt.Errorf("failed to force version: %w", err)
	}

	fmt.Printf("Database forced to version %d\n", version)
	return nil
}

func runMigrateVersion(cmd *cobra.Command, args []string) error {
	m, err := getMigrationInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		fmt.Println("No version set (database not initialized)")
		return nil
	}

	status := "clean"
	if dirty {
		status = "dirty"
	}

	fmt.Printf("Current version: %d (%s)\n", version, status)
	return nil
}
