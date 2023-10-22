package email

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(email string, subject string, body string) error {
	// Baca variabel lingkungan untuk konfigurasi SMTP
	smtpServer := os.Getenv("SMTPSERVER")
	smtpPortStr := os.Getenv("SMTPPORT")
	smtpUsername := os.Getenv("SMTPUSERNAME")
	smtpPassword := os.Getenv("SMTPPASSWORD")

	// Konversi smtpPortStr menjadi int
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUsername)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
