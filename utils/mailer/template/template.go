package template

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"sync"

	"github.com/ngs24313/gopu/config"
	"github.com/ngs24313/gopu/utils/log"
	"go.uber.org/zap"
)

//EmailTemplate email template
type EmailTemplate struct {
	Subject string
	Templ   *template.Template
}

//EmailTemplateResult result
type EmailTemplateResult struct {
	Subject string
	Body    string
}

var (
	templateMapping = map[string]*EmailTemplate{}
	bytBufPool      = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

//Init initialize template file
func Init(conf *config.Config) error {
	emailConf := conf.Mailer

	templates := emailConf.EmailTemplates.Templates

	for name, tmpl := range templates {
		emilTempl, err := loadEmailTemplate(&tmpl, conf)
		if err != nil {
			log.Logger(context.Background()).Warn("Failed to loading email template",
				zap.String("name", name),
				zap.Error(err))
			continue
		}
		templateMapping[name] = emilTempl
	}
	return nil
}

func loadEmailTemplate(tmpl *config.EmailTemplate, conf *config.Config) (*EmailTemplate, error) {
	var emailTempl EmailTemplate
	emailTempl.Subject = tmpl.Subject
	if tmpl.Filepath != "" {
		filepath := filepath.Join(conf.EmailTemplBasePath(), tmpl.Filepath)
		templ, err := template.ParseFiles(filepath)
		if err != nil {
			return nil, err
		}
		emailTempl.Templ = templ
	} else if tmpl.Content != "" {
		templ, err := template.New("").Parse(tmpl.Content)
		if err != nil {
			return nil, err
		}
		emailTempl.Templ = templ
	} else {
		return nil, fmt.Errorf("email template is not content")
	}
	return &emailTempl, nil
}

//GenEmailContent generate a html page for email
func GenEmailContent(name string, data interface{}) EmailTemplateResult {
	buff := bytBufPool.Get().(*bytes.Buffer)
	defer bytBufPool.Put(buff)
	buff.Reset()

	tmpl, ok := templateMapping[name]
	if !ok {
		log.Logger(context.Background()).Warn(
			"Cannot generate email content, name does not exists",
			zap.String("name", name),
			zap.Any("data", data),
		)
		return EmailTemplateResult{}
	}

	if err := tmpl.Templ.Execute(buff, data); err != nil {
		log.Logger(context.Background()).Warn("Failed to generate email content", zap.Error(err))
		return EmailTemplateResult{}
	}
	return EmailTemplateResult{
		Subject: tmpl.Subject,
		Body:    buff.String(),
	}
}
