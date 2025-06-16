package main

import (
	"log"

	"YoPost/internal/service"
)

func main() {
	to := "test@example.com"
	subject := "Test Subject"
	body := "Test Body"

	if err := service.SendTestEmail(to, subject, body); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}
	log.Println("Email sent successfully")
}
