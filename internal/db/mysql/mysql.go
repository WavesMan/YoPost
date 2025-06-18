package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type MySQLClient struct {
	db *sql.DB
}

func NewMySQLClient(config MySQLConfig) (*MySQLClient, error) {
	log.Printf("[INFO] Connecting to MySQL database at %s:%d", config.Host, config.Port)
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to open MySQL connection: %v", err)
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Printf("[ERROR] MySQL connection ping failed: %v", err)
		return nil, err
	}

	log.Printf("[INFO] Successfully connected to MySQL database")

	client := &MySQLClient{db: db}
	if err := client.initTables(); err != nil {
		log.Printf("[ERROR] Failed to initialize tables: %v", err)
		return nil, err
	}

	return client, nil
}

func (c *MySQLClient) initTables() error {
	log.Println("[DEBUG] Initializing MySQL tables")
	
	// 创建用户表
	_, err := c.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	log.Println("[INFO] MySQL tables initialized successfully")
	return nil
}

func (c *MySQLClient) GetDB() *sql.DB {
	return c.db
}

func (c *MySQLClient) Close() error {
	log.Println("[INFO] Closing MySQL connection")
	return c.db.Close()
}