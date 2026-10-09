package infrastructure

import (
	"github.com/velonyapp/email/internal/infrastructure/data/mysql"
	"github.com/velonyapp/email/internal/infrastructure/email/resend"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	resend.NewClient,
	resend.NewSender,
	mysql.NewConnection,
	mysql.NewUnitOfWork,
	mysql.NewTemplateRepo,
)
