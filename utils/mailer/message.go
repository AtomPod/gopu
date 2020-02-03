package mailer

import (
	"errors"

	"gopkg.in/gomail.v2"
)

var (
	//ErrMsgHasNoSender no sender
	ErrMsgHasNoSender = errors.New("message has no sender")
	//ErrMsgHasNoReceiver no receiver
	ErrMsgHasNoReceiver = errors.New("message has no receiver")
	//ErrMsgSubjectIsEmpty no subject
	ErrMsgSubjectIsEmpty = errors.New("message subject is empty")
	//ErrMsgBodyIsEmpty no body
	ErrMsgBodyIsEmpty = errors.New("message body is empty")
)

//User email user
type User struct {
	Name    string
	Address string
}

//Message email message
type Message struct {
	From User
	To   []User
	Cc   []User

	Subject     string
	Body        string
	ContentType string
}

//IsValid return true if message is valid, otherwise, return false
func (m *Message) IsValid() bool {
	if m.From.Address == "" || len(m.To) == 0 || m.Subject == "" || m.Body == "" {
		return false
	}
	return true
}

//ToGoMailMessage convert message to gomail.Message
func (m *Message) ToGoMailMessage() (*gomail.Message, error) {
	if m.From.Address == "" {
		return nil, ErrMsgHasNoSender
	}

	if len(m.Subject) == 0 {
		return nil, ErrMsgSubjectIsEmpty
	}

	if len(m.Body) == 0 {
		return nil, ErrMsgBodyIsEmpty
	}

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", m.From.Address, m.From.Name)

	to := make([]string, 0)
	for _, t := range m.To {
		if t.Address != "" {
			to = append(to, t.Address)
		}
	}

	if len(to) == 0 {
		return nil, ErrMsgHasNoReceiver
	}
	msg.SetHeader("To", to...)

	cc := make([]string, 0)
	for _, c := range m.Cc {
		if c.Address != "" {
			cc = append(cc, c.Address)
		}
	}

	if len(cc) > 0 {
		msg.SetHeader("Cc", cc...)
	}

	msg.SetHeader("Subject", m.Subject)
	msg.SetBody(m.ContentType, m.Body)

	return msg, nil
}
