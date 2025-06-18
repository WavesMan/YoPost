package main

import (
	"log"
	"net/http"

	"YoPost/internal/api/smtp"
	"YoPost/internal/db"
	"YoPost/internal/db/mongodb"
	"YoPost/internal/db/mysql"
)

func main() {
	// Initialize database
	dbConfig := db.DBConfig{
		MySQL: mysql.MySQLConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "password",
			Database: "yopost",
		},
		MongoDB: mongodb.MongoDBConfig{
			Host:     "localhost",
			Port:     27017,
			User:     "admin",
			Password: "password",
			Database: "yopost",
		},
	}
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
