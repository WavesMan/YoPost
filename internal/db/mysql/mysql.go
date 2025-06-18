package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-sql-driver/mysql"
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
	const maxRetries = 3
	const retryDelay = 5 * time.Second

	// 增强密码验证逻辑
	if len(config.Password) < 8 {
		log.Printf("[WARN] MySQL password is too short (minimum 8 characters)")
	}

	// Add config validation
	if config.Host == "" {
		return nil, fmt.Errorf("MySQL host cannot be empty")
	}
	if config.User == "" || config.Password == "" {
		return nil, fmt.Errorf("MySQL credentials cannot be empty")
	}
	if config.Database == "" {
		return nil, fmt.Errorf("MySQL database name cannot be empty")
	}

	log.Printf("[DEBUG] Using MySQL config: Host=%s, Port=%d, User=%s, Database=%s", 
		config.Host, config.Port, config.User, config.Database)

	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		log.Printf("[INFO] Attempting MySQL connection (try %d/%d) to %s:%d", 
			i+1, maxRetries, config.Host, config.Port)

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&timeout=5s",
			config.User, config.Password, config.Host, config.Port, config.Database)

		// Add debug log for actual connection string (mask password)
		log.Printf("[DEBUG] Using DSN: %s:***@tcp(%s:%d)/%s", 
			config.User, config.Host, config.Port, config.Database)

		// Add network diagnostics
		log.Printf("[DEBUG] Network configuration before connection attempt")
		if addrs, err := net.LookupHost(config.Host); err == nil {
			log.Printf("[DEBUG] Resolved host %s to: %v", config.Host, addrs)
		} else {
			log.Printf("[DEBUG] Failed to resolve host %s: %v", config.Host, err)
		}

		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("[ERROR] MySQL connection failed: %v", err)
			time.Sleep(retryDelay)
			continue
		}

		// 优化连接池配置
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(10)
		db.SetConnMaxLifetime(30 * time.Minute)
		db.SetConnMaxIdleTime(10 * time.Minute)

		if err = db.Ping(); err != nil {
			log.Printf("[ERROR] MySQL connection ping failed: %v", err)
			// 增强错误提示
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				log.Printf("[ERROR] MySQL authentication failed for user '%s' from %s", 
					config.User, config.Host)
				log.Printf("[DEBUG] MySQL error details: Number=%d, SQLState=%s, Message=%s", 
					mysqlErr.Number, mysqlErr.SQLState, mysqlErr.Message)
			}
			db.Close()
			time.Sleep(retryDelay)
			continue
		}

		log.Printf("[INFO] Successfully connected to MySQL database")
		break
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL after %d attempts: %v", maxRetries, err)
	}

	client := &MySQLClient{db: db}
	if err := client.initTables(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize tables: %v", err)
	}

	return client, nil
}

// 添加initTables方法实现
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
