package service

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

//
const (
	TextHTML  = "text/html"
	TextPlain = "text/plain"
)

// SendWithoutSSL send mail
func SendWithoutSSL(smtpAddress string, port int, username, password string, message *Message) error {
	return send(smtpAddress, port, username, password, false, message)
}

func send(smtpAddress string, port int, username, password string, ssl bool, message *Message) error {
	if err := verify(message); err != nil {
		return err
	}
	var m = gomail.NewMessage()
	setContact(m, "From", message.From)
	setContact(m, "To", message.To...)
	setContact(m, "Cc", message.Cc...)
	setContact(m, "Bcc", message.Bcc...)
	m.SetHeader("Subject", message.Subject)
	m.SetBody(message.ContentType, message.Body)
	setFile(m, message.Attach)

	var dialer = gomail.Dialer{
		Host:     smtpAddress,
		Port:     port,
		Username: username,
		Password: password,
		SSL:      ssl,
	}
	return dialer.DialAndSend(m)
}

func setFile(m *gomail.Message, files []File) {
	if len(files) == 0 {
		return
	}
	for _, file := range files {
		if file.Path == "" {
			continue
		}
		if file.AliasName == "" {
			m.Attach(file.Path)
		} else {
			m.Attach(file.Path, gomail.Rename(file.AliasName))
		}
	}
}

func setContact(m *gomail.Message, field string, contacts ...Contact) {
	if len(contacts) == 0 {
		return
	}
	var values = make([]string, 0, len(contacts))
	for _, contact := range contacts {
		if contact.Address == "" {
			continue
		}
		if contact.Name == "" {
			values = append(values, contact.Address)
		} else {
			values = append(values, m.FormatAddress(contact.Address, contact.Name))
		}
	}
	if len(values) == 0 {
		return
	}
	m.SetHeader(field, values...)
}

func verify(message *Message) error {
	if message.From.Address == "" {
		return fmt.Errorf("from address is missing")
	}

	var hasTo bool
	for _, to := range message.To {
		if to.Address != "" {
			hasTo = true
		}
	}
	if !hasTo {
		return fmt.Errorf("to address is missing")
	}

	switch message.ContentType {
	case TextHTML, TextPlain:
	default:
		return fmt.Errorf("ContentType is invalid")
	}
	return nil
}

// Message msg
type Message struct {
	From        Contact
	To          []Contact
	Cc          []Contact
	Bcc         []Contact
	Subject     string
	Body        string
	ContentType string
	Attach      []File
}

// Contact contact
type Contact struct {
	Name    string
	Address string
}

// File file
type File struct {
	Path      string
	AliasName string
}
