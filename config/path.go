package config

import (
	"os"
	"path/filepath"
)

//InitBaseDir create all directory if not exists
func (c *Config) InitBaseDir() error {
	if err := os.MkdirAll(c.AvatarBasePath(), 0666); err != nil {
		return err
	}
	if err := os.MkdirAll(c.EmailTemplBasePath(), 0666); err != nil {
		return err
	}
	return nil
}

//AvatarBasePath for avatar file path
func (c *Config) AvatarBasePath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "/images/avatar")
}

//EmailTemplBasePath for email template file path
func (c *Config) EmailTemplBasePath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "/templates/email")
}
