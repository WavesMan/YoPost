package mail

import (
	"encoding/json"
	"net/http"
)

type SMTPRequest struct {
	From       string   `json:"from"`
	To         []string `json:"to"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	Extensions struct {
		STARTTLS bool `json:"starttls"`
		EightBit bool `json:"8bitmime"`
		SMTPUTF8 bool `json:"smtputf8"`
	} `json:"extensions"`
}

type SMTPResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	MessageID string `json:"message_id,omitempty"`
	Error     string `json:"error,omitempty"`
}

type MailServer struct {
	core Core
}

func NewMailServer(core Core) *MailServer {
	return &MailServer{core: core}
}

func (s *MailServer) HandleSMTPSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ... existing request processing code from handler.go ...
}

func sendErrorResponse(w http.ResponseWriter, errorCode, message string, status int) {
	resp := SMTPResponse{
		Success: false,
		Error:   errorCode,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
