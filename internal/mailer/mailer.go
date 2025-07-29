package mailer

import (
	"fmt"
	"net/smtp"
)

type SMTPMailer struct {
	From     string
	Password string
	Host     string
	Port     string
}

func NewSMTPMailer(from, password, host, port string) *SMTPMailer {
	return &SMTPMailer{
		From:     from,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (m *SMTPMailer) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.From, m.Password, m.Host)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	return smtp.SendMail(addr, auth, m.From, []string{to}, msg)
}
