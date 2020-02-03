package mailer

import (
	"context"
	"crypto/tls"

	"github.com/ngs24313/gopu/config"
	"github.com/ngs24313/gopu/utils/log"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type goMailer struct {
	conf config.Mailer
}

//NewGOMailer create gomail mailer
func NewGOMailer(conf config.Config) Mailer {
	return &goMailer{
		conf: conf.Mailer,
	}
}

func (m *goMailer) dialer() *gomail.Dialer {
	conf := m.conf
	dialer := gomail.NewDialer(conf.Host, conf.Port, conf.Username, conf.Password)

	if conf.TLS != nil {
		cert, err := tls.LoadX509KeyPair(conf.TLS.CertPath, conf.TLS.KeyPath)
		if err == nil {
			dialer.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
			dialer.SSL = true
		} else {
			log.Logger(context.Background()).Warn("Failed to use ssl for email", zap.Error(err))
		}
	}

	return dialer
}

func (m *goMailer) Send(msg *Message) error {
	message, err := msg.ToGoMailMessage()
	if err != nil {
		return err
	}
	return m.dialer().DialAndSend(message)
}
