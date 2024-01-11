package mailer

import (
	"bytes"
	"fmt"
	"github.com/lambda-platform/lambda/config"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func NewRequest(to []string, subject string, from string) *Request {

	return &Request{
		to:      to,
		from:    from,
		subject: subject,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail() bool {

	fromEmail := r.from
	toEmail := r.to[0]
	smtpServer := config.Config.Mail.Host
	smtpPort := config.Config.Mail.Port
	smtpUser := config.Config.Mail.Username
	smtpPassword := config.Config.Mail.Password

	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", r.subject)
	m.SetBody("text/html", r.body)

	// Send the email using an SMTP server
	d := gomail.NewDialer(smtpServer, smtpPort, smtpUser, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
		return false
	} else {
		return true
	}

}

func (r *Request) Send(templateName string, items interface{}) bool {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	if ok := r.sendMail(); ok {
		return true
	} else {
		return false
	}
}
func (r *Request) SendByTemplate(templateString string) bool {

	r.body = templateString
	if ok := r.sendMail(); ok {
		return true
	} else {
		return false
	}
}
