package main

import (
	"log"
	"net/http"

	"YoPost/internal/api/smtp"
	"YoPost/internal/db"
	"YoPost/internal/db/mongodb"
	"YoPost/internal/db/mysql"
)

var (
	dbConfig = db.DBConfig{
		MySQL: mysql.MySQLConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			User:     "yopost",
			Password: "123456",
			Database: "yopost",
		},
		MongoDB: mongodb.MongoDBConfig{
			Host:     "127.0.0.1",
			Port:     27017,
			User:     "yopost",
			Password: "123456",
			Database: "yopost",
		},
	}
)

func main() {
	// Initialize database
	db.InitDB(dbConfig)

	// Initialize API routes
	http.HandleFunc("/api/smtp/send", smtp.SendEmailHandler)
	http.HandleFunc("/api/smtp/config", smtp.GetConfigHandler)

	// Start API server
	log.Println("Starting API server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
