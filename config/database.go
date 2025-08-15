package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DatabaseConfig 數據庫配置結構
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDatabaseConfig 創建新的數據庫配置
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "123456"),
		DBName:   getEnv("DB_NAME", "fetu_bo"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// Connect 連接到數據庫
func (c *DatabaseConfig) Connect() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 測試連接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

// getEnv 獲取環境變量，如果不存在則返回默認值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
