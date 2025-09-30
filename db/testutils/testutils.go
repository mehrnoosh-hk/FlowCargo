package testutils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	once      sync.Once
	dbManager *TestDBManager
	initError error
)

// TestDB interface defines methods for interacting with a test database
type TestDB interface {
	GetPool() *pgxpool.Pool
	BeginTx(ctx context.Context, t *testing.T) pgx.Tx
	Close()
}

// TestDBManager implements the TestDB interface
// encapsulates all database management logic
type TestDBManager struct {
	// Private fields - encapsulated state
	dsn         string
	db          *pgxpool.Pool
	initialized bool
	closed      bool
}

func GetDBManager() TestDB {
	once.Do(func() {
		dbManager = NewTestDBManager(GetTestDSN())
		initError = dbManager.Initialize()
	})
	return dbManager
}

func NewTestDBManager(dsn string) *TestDBManager {
	return &TestDBManager{
		dsn:         dsn,
		db:          nil,
		initialized: false,
		closed:      false,
	}
}

// Initialize sets up the database connection (called once)
func (m *TestDBManager) Initialize() error {
	var initErr error
	ctx := context.Background()

	if m.closed {
		initErr = fmt.Errorf("database manager has been closed")
		return initErr
	}

	// Open database connection
	db, err := pgxpool.New(ctx, m.dsn)
	if err != nil {
		initErr = fmt.Errorf("failed to create database connection: %w", err)
		return initErr
	}

	// Test connection
	if err := db.Ping(ctx); err != nil {
		db.Close()
		initErr = fmt.Errorf("failed to ping database: %w", err)
		return initErr
	}

	// Run migrations
	if err := migrateDatabase(m.dsn, true); err != nil {
		db.Close()
		initErr = fmt.Errorf("failed to run migrations: %w", err)
		return initErr
	}

	m.db = db
	m.initialized = true

	fmt.Println("Test database initialized")

	return initErr
}

func (m *TestDBManager) BeginTx(ctx context.Context, t *testing.T) pgx.Tx {
	if t != nil {
		t.Helper()
	}

	db := m.getDB(t)

	tx, err := db.Begin(ctx)
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}
		panic(err)
	}

	return tx
}

func (m *TestDBManager) getDB(t *testing.T) *pgxpool.Pool {
	if t != nil {
		t.Helper()
	}
	if !m.initialized {
		if err := m.Initialize(); err != nil {
			if t != nil {
				t.Fatalf("Failed to initialize database: %v", err)
			}
			panic(err)
		}
	}

	if m.closed {
		if t != nil {
			t.Fatal("Database manager has been closed")
		}
		panic("Database manager has been closed")
	}

	if !m.initialized || m.db == nil {
		if t != nil {
			t.Fatal("Database not properly initialized")
		}
		panic("Database not properly initialized")
	}

	return m.db
}

// Close shuts down the database connection
func (m *TestDBManager) Close() {

	// Run down migrations
	err := migrateDatabase(m.dsn, false) // TODO: make sure wether it is required or not
	if err != nil {
		fmt.Printf("Failed to run down migrations: %v\n", err)
	}

	if m.db != nil {
		fmt.Println("Closing test database connection...")
		m.db.Close()
		m.db = nil
	}

	m.closed = true
	fmt.Println("Test Resources has been cleaned successfully!")
}

// GetTestDSN returns the database DSN for tests from environment variables
func GetTestDSN() string {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		// Default test database URL - matches Docker setup
		dsn = "postgres://postgres:password@localhost:5432/flowcargo_test?sslmode=disable&search_path=public"
	}
	return dsn
}

func migrateDatabase(URL string, up bool) error {
	migrationsPath, err := findMigrationsDir()
	if err != nil {
		return fmt.Errorf("failed to find migrations directory: %w", err)
	}
	migrationsURL := "file://" + migrationsPath

	m, err := migrate.New(migrationsURL, URL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if up {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		return nil
	}

	if err := m.Drop(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

// findMigrationsDir searches for the migrations directory by walking up the directory tree
// It looks for a directory containing both go.mod (project root) and migrations folder
func findMigrationsDir() (string, error) {
	// Get the directory of this source file
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("failed to get caller information")
	}

	currentDir := filepath.Dir(filename)

	// Walk up the directory tree
	for {
		// Check if go.mod exists in current directory (indicates project root)
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			// Found project root, now check for migrations directory
			migrationsPath := filepath.Join(currentDir, "migrations")
			if _, err := os.Stat(migrationsPath); err == nil {
				return migrationsPath, nil
			}
			return "", fmt.Errorf("go.mod found at %s but migrations directory not found", currentDir)
		}

		// Move up one directory
		parentDir := filepath.Dir(currentDir)

		// Check if we've reached the filesystem root
		if parentDir == currentDir {
			return "", fmt.Errorf("reached filesystem root without finding go.mod and migrations directory")
		}

		currentDir = parentDir
	}
}

func (m *TestDBManager) GetPool() *pgxpool.Pool {
	return m.db
}
