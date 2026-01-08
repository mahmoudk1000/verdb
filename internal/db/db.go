package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"

	"github.com/mahmoudk1000/relen/internal/database"
)

var (
	instance *database.Queries
	conn     *sql.DB
	once     sync.Once
	initErr  error
)

func Init(connectionString string) error {
	once.Do(func() {
		conn, initErr = sql.Open("postgres", connectionString)
		if initErr != nil {
			initErr = fmt.Errorf("failed to open database: %w", initErr)
			return
		}

		if err := conn.Ping(); err != nil {
			initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		instance = database.New(conn)
	})

	return initErr
}

func Get() *database.Queries {
	if instance == nil {
		panic("database not initialized: call db.Init() first")
	}
	return instance
}

func GetConn() *sql.DB {
	if conn == nil {
		panic("database not initialized: call db.Init() first")
	}
	return conn
}

func Close() error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}
