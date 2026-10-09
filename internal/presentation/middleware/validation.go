package middleware

import (
	"errors"

	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/validate"
	"github.com/velonyapp/email/gen/api/v1"
	"go.einride.tech/aip/fieldbehavior"
	"google.golang.org/protobuf/proto"
)

type Validation middleware.Middleware

func NewValidation() Validation {
	return Validation(
		validate.Validator(func(req any) error {
			switch req := req.(type) {
			case *apiv1.UpdateEmailTemplateRequest:
				if req.GetEmailTemplate() == nil {
					return errors.New("missing required field: email template")
				}

				return fieldbehavior.ValidateRequiredFieldsWithMask(
					req.GetEmailTemplate(),
					req.GetUpdateMask(),
				)

			default:
				message, ok := req.(proto.Message)
				if !ok {
					return nil
				}

				return fieldbehavior.ValidateRequiredFields(message)
			}
		}),
	)
}

func (m Validation) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
