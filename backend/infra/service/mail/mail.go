package mail

import (
	"fmt"

	"github.com/eclipse-disuko/disuko/domain/mailtemplate"
	"github.com/eclipse-disuko/disuko/helper/mail"
	"github.com/eclipse-disuko/disuko/infra/repository/mailtemplates"
	"github.com/eclipse-disuko/disuko/logy"
)

type Service struct {
	Client       *mail.Client
	TemplateRepo mailtemplates.IMailTemplatesRepository
}

func (s *Service) SendMail(rs *logy.RequestSession, to string, templateKey mailtemplate.MailTemplateKey, data any) error {
	tmpl := s.TemplateRepo.FindByKey(rs, string(templateKey), false)
	if tmpl == nil {
		return fmt.Errorf("template not found %s", templateKey)
	}
	if err := s.Client.Send(to, tmpl.Cc, tmpl.Bcc, tmpl.Message, tmpl.Subject, data); err != nil {
		return fmt.Errorf("sending mail: %w", err)
	}
	return nil
}

func (s *Service) SendTestMail(rs *logy.RequestSession, to string, templateKey mailtemplate.MailTemplateKey, message string) error {
	tmpl := s.TemplateRepo.FindByKey(rs, string(templateKey), false)
	if tmpl == nil {
		return fmt.Errorf("template not found %s", templateKey)
	}
	data := make(map[string]string, len(tmpl.Values))
	for key, description := range tmpl.Values {
		data[key] = fmt.Sprintf("<< %s >>", description)
	}
	if err := s.Client.Send(to, "", "", message, tmpl.Subject, data); err != nil {
		return fmt.Errorf("sending test mail: %w", err)
	}
	return nil
}
