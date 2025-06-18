package db

import (
	"log"
	
	"YoPost/internal/db/mongodb"
	"YoPost/internal/db/mysql"
)

type DBConfig struct {
	MySQL    mysql.MySQLConfig
	MongoDB  mongodb.MongoDBConfig
}

func InitDB(config DBConfig) {
	// Initialize MySQL
	mysqlClient, err := mysql.NewMySQLClient(config.MySQL)
	if err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}
	defer mysqlClient.Close()

	// Initialize MongoDB
	mongoClient, err := mongodb.NewMongoDBClient(config.MongoDB)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer mongoClient.Close()

	log.Println("Database services initialized successfully")
}