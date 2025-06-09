package protocol

import (
	"fmt"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
)

// MockMailCore 实现mail.Core接口，用于测试
type MockMailCore struct {
	StoredEmails []struct {
		From       string
		Recipients []string
		Data       string
	}
}

func (m *MockMailCore) ValidateUser(email string) bool { return true }

func (m *MockMailCore) StoreEmail(From string, Recipients []string, Data string) error {
	m.StoredEmails = append(m.StoredEmails, struct {
		From       string
		Recipients []string
		Data       string
	}{From, Recipients, Data})
	return nil
}

func (m *MockMailCore) GetConfig() *config.Config { return &config.Config{} }

func (m *MockMailCore) GetEmail(id string) (*mail.Email, error) {
	if len(m.StoredEmails) > 0 {
		return &mail.Email{
			ID:   id,
			Body: m.StoredEmails[0].Data,
			From: m.StoredEmails[0].From,
			To:   m.StoredEmails[0].Recipients,
		}, nil
	}
	return nil, fmt.Errorf("email not found")
}

func (m *MockMailCore) GetEmails() ([]mail.Email, error) {
	emails := make([]mail.Email, len(m.StoredEmails))
	for i, stored := range m.StoredEmails {
		emails[i] = mail.Email{
			ID:   fmt.Sprintf("email-%d", i),
			Body: stored.Data,
			From: stored.From,
			To:   stored.Recipients,
		}
	}
	return emails, nil
}
