package config

import (
	"crypto/tls"
	"fmt"

	"golang.org/x/crypto/acme/autocert"
)

//GenTLSConfig generate TLS Config
func (c *Config) GenTLSConfig() (*tls.Config, error) {
	conf := c
	tlsConf := conf.HTTP.TLS

	if tlsConf.CertPath != "" && tlsConf.KeyPath != "" {
		cert, err := tls.LoadX509KeyPair(tlsConf.CertPath, tlsConf.KeyPath)
		if err != nil {
			return nil, err
		}
		return &tls.Config{
			Certificates: []tls.Certificate{cert},
		}, nil
	}

	acme := conf.HTTP.Acme
	if acme.Type != "" {
		if acme.Type == "Let's Encrypt" {
			m := autocert.Manager{
				Cache:      autocert.DirCache(acme.DirCache),
				HostPolicy: autocert.HostWhitelist(acme.Hosts...),
				Prompt:     autocert.AcceptTOS,
			}
			return m.TLSConfig(), nil
		}
		return nil, fmt.Errorf("Unknow ACME Type: %s", acme.Type)
	}
	return nil, nil
}
