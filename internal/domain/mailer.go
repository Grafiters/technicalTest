package domain

type MailerInterface interface {
	SendEmail(to string, subject string, body string) error
}
