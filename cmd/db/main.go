package main

import (
	"log"
	"time"
	
	"YoPost/internal/db/mongodb"
	"YoPost/internal/db/mysql"
)

type DBConfig struct {
	MySQL    mysql.MySQLConfig
	MongoDB  mongodb.MongoDBConfig
}

func InitDB(config DBConfig) {
	const maxRetries = 3
	const retryInterval = 5 * time.Second

	// Initialize MySQL with retry logic
	var mysqlClient *mysql.MySQLClient
	var err error
	for i := 0; i < maxRetries; i++ {
		mysqlClient, err = mysql.NewMySQLClient(config.MySQL)
		if err == nil {
			break
		}
		log.Printf("MySQL connection attempt %d failed: %v", i+1, err)
		time.Sleep(retryInterval)
	}
	if err != nil {
		log.Fatalf("Failed to initialize MySQL after %d attempts: %v", maxRetries, err)
	}
	defer mysqlClient.Close()

	// Initialize MongoDB with retry logic
	var mongoClient *mongodb.MongoDBClient
	for i := 0; i < maxRetries; i++ {
		mongoClient, err = mongodb.NewMongoDBClient(config.MongoDB)
		if err == nil {
			break
		}
		log.Printf("MongoDB connection attempt %d failed: %v", i+1, err)
		time.Sleep(retryInterval)
	}
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB after %d attempts: %v", maxRetries, err)
	}
	defer mongoClient.Close()

	log.Println("Database services initialized successfully")
}