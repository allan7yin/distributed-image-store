package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectionHandler manages the database connection pool
type ConnectionHandler struct {
	DB   *gorm.DB
	Pool *pgxpool.Pool
}

// NewConnectionHandler sets up the connection pool and GORM instance
func NewConnectionHandler() (*ConnectionHandler, error) {
	// Load properties from environment variables or a config file
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("database URL is not set")
	}

	// Set up the connection pool using pgxpool
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
		return nil, err
	}

	// Set connection pool options
	poolConfig.MaxConns = 20 // max pool size
	poolConfig.MinConns = 5  // min idle connections
	poolConfig.MaxConnIdleTime = 30 * time.Second
	poolConfig.MaxConnLifetime = 60 * time.Minute
	poolConfig.HealthCheckPeriod = 30 * time.Second

	// Establish the connection pool
	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to establish database connection pool: %v", err)
		return nil, err
	}

	// Use pgxpool to create a database/sql connection (for GORM)
	sqlDB := stdlib.OpenDB(*poolConfig.ConnConfig)

	// Set up GORM using the *sql.DB instance from pgxpool
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Show SQL
	})
	if err != nil {
		log.Fatalf("Error setting up GORM: %v", err)
		return nil, err
	}

	return &ConnectionHandler{
		DB:   gormDB,
		Pool: pool,
	}, nil
}

// OpenTransaction opens a new GORM session with a transaction
func (handler *ConnectionHandler) OpenTransaction() (*gorm.DB, func() error, func() error, error) {
	tx := handler.DB.Begin()
	if tx.Error != nil {
		return nil, nil, nil, tx.Error
	}

	commit := func() error {
		return tx.Commit().Error
	}

	rollback := func() error {
		return tx.Rollback().Error
	}

	return tx, commit, rollback, nil
}

// Close closes the database connection pool
func (handler *ConnectionHandler) Close() {
	handler.Pool.Close()
}

// PrintConnectionPoolStats prints the stats of the connection pool
func (handler *ConnectionHandler) PrintConnectionPoolStats() {
	stats := handler.Pool.Stat()
	fmt.Printf("Total Connections: %d\n", stats.TotalConns())
	fmt.Printf("Idle Connections: %d\n", stats.IdleConns())
	fmt.Printf("Active Connections: %d\n", stats.AcquiredConns())
	fmt.Printf("Threads Awaiting Connection: %d\n", stats.AcquireCount()) // Similar to threads awaiting
}
