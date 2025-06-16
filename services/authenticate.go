package service

import (
	"log"
	"net/smtp"
)

// Authenticate performs SMTP authentication with given credentials
// c: SMTP client connection
// host: SMTP server host
// username: authentication username
// password: authentication password
func Authenticate(c *smtp.Client, host, username, password string) error {
	if username != "" && password != "" {
		log.Printf("INFO: Attempting authentication for %s", username)
		auth := smtp.PlainAuth("", username, password, host)
		if err := c.Auth(auth); err != nil {
			log.Printf("ERROR: Authentication failed for %s - %v", username, err)
			return err
		}
		log.Printf("INFO: Successfully authenticated %s", username)
	}
	return nil
}