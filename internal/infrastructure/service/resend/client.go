package resend

import (
	"github.com/velonyapp/email/internal/conf"

	"github.com/resend/resend-go/v3"
)

func NewClient(c *conf.Service) *resend.Client {
	return resend.NewClient(c.Resend.ApiKey)
}
