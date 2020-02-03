package mailer

import "github.com/ngs24313/gopu/config"

var (
	defaultMailer Mailer
)

//Mailer mailer interface
type Mailer interface {
	Send(msg *Message) error
}

//Init initialize mailer
func Init(conf *config.Config) error {
	defaultMailer = NewGOMailer(*conf)
	return nil
}

//Send send the message by default mailer
func Send(msg *Message) error {
	return defaultMailer.Send(msg)
}
