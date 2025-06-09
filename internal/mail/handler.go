package mail

import (
	"net/http"
)

type MailHandler struct {
	mailServer *MailServer
}

func NewMailHandler(mailCore Core) *MailHandler {
	return &MailHandler{
		mailServer: NewMailServer(mailCore),
	}
}

func (h *MailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/smtp/send":
		h.mailServer.HandleSMTPSend(w, r)
	default:
		http.NotFound(w, r)
	}
}