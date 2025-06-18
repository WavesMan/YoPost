package smtp

import (
	"YoPost/internal/mail/core"
	"encoding/json"
	"log"
	"net/http"
)

// SendEmailRequest defines the request structure for sending emails
type SendEmailRequest struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// SendEmailResponse defines the response structure for sending emails
type SendEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ConfigResponse defines the response structure for SMTP configuration
type ConfigResponse struct {
	Host      string `json:"host"`
	TLSPort   string `json:"tls_port"`
	NoTLSPort string `json:"no_tls_port"`
}

// SendEmailHandler handles email sending requests
func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	log.Printf("INFO: Handling email send request")
	var req SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body - %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("DEBUG: Request details - To: %s, Subject: %s", req.To, req.Subject)

	// Build message
	msg := []byte("To: " + req.To + "\r\n" +
		"Subject: " + req.Subject + "\r\n" +
		"\r\n" +
		req.Body + "\r\n")
	log.Printf("DEBUG: Built message (%d bytes)", len(msg))

	// Send email
	log.Printf("INFO: Attempting to send email to %s", req.To)
	err := core.TLSstatus(req.Username, []string{req.To}, msg, req.Username, req.Password)
	if err != nil {
		log.Printf("ERROR: Failed to send email to %s - %v", req.To, err)
		response := SendEmailResponse{
			Success: false,
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	log.Printf("INFO: Successfully sent email to %s", req.To)
	response := SendEmailResponse{
		Success: true,
		Message: "Email sent successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetConfigHandler returns SMTP server configuration
func GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: Handling SMTP config request")
	config := core.GetMailServerConfig()
	if config == nil {
		log.Printf("ERROR: SMTP configuration not initialized")
		http.Error(w, "SMTP configuration not initialized", http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Config details - Host: %s, TLS Port: %s, NoTLS Port: %s", 
		config.Host, config.TLSPort, config.NoTLSPort)

	response := ConfigResponse{
		Host:      config.Host,
		TLSPort:   config.TLSPort,
		NoTLSPort: config.NoTLSPort,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
