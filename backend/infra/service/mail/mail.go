package mail

import (
	"fmt"
	"reflect"

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
	if expected, ok := mailtemplate.MailTemplateDataTypes[templateKey]; ok {
		if reflect.TypeOf(data) != reflect.TypeOf(expected) {
			return fmt.Errorf("mail data type %T does not match expected type %T for template %s", data, expected, templateKey)
		}
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
	var data any = map[string]string{}
	if expected, ok := mailtemplate.MailTemplateDataTypes[templateKey]; ok {
		data = FillWithPlaceholders(expected)
	}
	if err := s.Client.Send(to, "", "", message, tmpl.Subject, data); err != nil {
		return fmt.Errorf("sending test mail: %w", err)
	}
	return nil
}

func ValuesFor(templateKey mailtemplate.MailTemplateKey) map[string]string {
	expected, ok := mailtemplate.MailTemplateDataTypes[templateKey]
	if !ok {
		return map[string]string{}
	}
	return FieldDescriptions(expected)
}

func ToDto(entity *mailtemplate.MailTemplate) *mailtemplate.MailTemplateDto {
	dto := entity.ToDto()
	dto.Values = ValuesFor(mailtemplate.MailTemplateKey(entity.Key))
	return dto
}

func ToDtos(entities []*mailtemplate.MailTemplate) []*mailtemplate.MailTemplateDto {
	dtos := make([]*mailtemplate.MailTemplateDto, len(entities))
	for i, entity := range entities {
		dtos[i] = ToDto(entity)
	}
	return dtos
}
