package api

import (
	"log"
	"net/http"

	"YoPost/internal/api/smtp"
)

type APIConfig struct {
	Port string
}

func InitAPI(config APIConfig) {
	// Initialize API routes
	http.HandleFunc("/api/smtp/send", smtp.SendEmailHandler)
	http.HandleFunc("/api/smtp/config", smtp.GetConfigHandler)

	// Start API server
	log.Printf("Starting API server on :%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}