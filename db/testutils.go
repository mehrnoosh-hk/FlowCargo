package testutils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestDBManager encapsulates all database management logic
type TestDBManager struct {
	// Private fields - encapsulated state
	db          *pgxpool.Pool
	dsn         string
	initialized bool
	closed      bool
}

func NewTestDBManager(dsn string) *TestDBManager {
	return &TestDBManager{
		db:          nil,
		dsn:         dsn,
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

// GetDB returns the database connection, initializing if necessary
func (m *TestDBManager) GetDB(t *testing.T) *pgxpool.Pool {
	if t != nil {
		t.Helper()
	}

	// Initialize if needed
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
	// migrateDatabase(m.dsn, false) // TODO: make sure wether it is required or not

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
	// Get the directory of this source file
	_, filename, _, _ := runtime.Caller(0)
	sourceDir := filepath.Dir(filename)
	// Go up one level to project root, then to migrations
	migrationsPath := filepath.Join(sourceDir, "..", "migrations")
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
