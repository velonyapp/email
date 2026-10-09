package middleware

import (
	"context"
	"errors"

	applicationcommon "github.com/velonyapp/email/internal/application/common"
	domainentity "github.com/velonyapp/email/internal/domain/entity"
	domainservice "github.com/velonyapp/email/internal/domain/service"
	domainvo "github.com/velonyapp/email/internal/domain/vo"
	presentationapi "github.com/velonyapp/email/internal/presentation/api"

	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/middleware"
)

type ErrorMapper middleware.Middleware

func NewErrorMapper() ErrorMapper {
	return ErrorMapper(
		func(next middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req any) (any, error) {
				reply, err := next(ctx, req)
				if err == nil {
					return reply, nil
				}

				var kratosErr *kerrors.Error
				if errors.As(err, &kratosErr) {
					return reply, err
				}

				switch {
				// Template ID
				case errors.Is(err, domainvo.ErrTemplateIDInvalid):
					return reply, kerrors.BadRequest("",
						domainvo.ErrTemplateIDInvalid.Error(),
					)

				// Alias
				case errors.Is(err, domainvo.ErrAliasEmpty):
					return reply, kerrors.BadRequest("",
						domainvo.ErrAliasEmpty.Error(),
					)
				case errors.Is(err, domainvo.ErrAliasTooLong):
					return reply, kerrors.BadRequest("",
						domainvo.ErrAliasTooLong.Error(),
					)
				case errors.Is(err, domainvo.ErrAliasInvalidCharacter):
					return reply, kerrors.BadRequest("",
						domainvo.ErrAliasInvalidCharacter.Error(),
					)

				// Address
				case errors.Is(err, domainvo.ErrAddressInvalid):
					return reply, kerrors.BadRequest("",
						domainvo.ErrAddressInvalid.Error(),
					)

				// Template
				case errors.Is(err, domainentity.ErrTemplateBodyRequired):
					return reply, kerrors.BadRequest("",
						domainentity.ErrTemplateBodyRequired.Error(),
					)

				// Alias Policy
				case errors.Is(err, domainservice.ErrAliasAlreadyExists):
					return reply, kerrors.Conflict("",
						domainservice.ErrAliasAlreadyExists.Error(),
					)

				// Application Common
				case errors.Is(err, applicationcommon.ErrTemplateNotFound):
					return reply, kerrors.NotFound("",
						applicationcommon.ErrTemplateNotFound.Error(),
					)

				// API Service
				case errors.Is(err, presentationapi.ErrInvalidEmailTemplateResourceName):
					return reply, kerrors.BadRequest("",
						presentationapi.ErrInvalidEmailTemplateResourceName.Error(),
					)

				case errors.Is(err, presentationapi.ErrInvalidUpdateMaskPath):
					return reply, kerrors.BadRequest("",
						presentationapi.ErrInvalidUpdateMaskPath.Error(),
					)

				default:
					return reply, kerrors.InternalServer("",
						"internal server error",
					).WithCause(err)
				}
			}
		},
	)
}

func (m ErrorMapper) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
