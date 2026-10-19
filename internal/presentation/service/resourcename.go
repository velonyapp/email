
package service

import (
	"errors"

	"go.einride.tech/aip/resourcename"
)

var (
	ErrInvalidEmailTemplateResourceName = errors.New("invalid email template resource name")
)

const (
	emailTemplateResourcePattern = "emailTemplates/{email_template}"
)

func parseEmailTemplateName(name string) (string, error) {
	var templateID string
	if err := resourcename.Sscan(name, emailTemplateResourcePattern, &templateID); err != nil {
		return "", ErrInvalidEmailTemplateResourceName
	}

	return templateID, nil
}

func formatEmailTemplateName(templateID string) string {
	return resourcename.Sprint(emailTemplateResourcePattern, templateID)
}
