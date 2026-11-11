package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/velonyapp/email/internal/domain/entity"
	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"

	"github.com/go-sql-driver/mysql"
)

var _ repo.Template = (*templateRepo)(nil)

type templateRepo struct {
	db  *sql.DB
	log *slog.Logger
}

func NewTemplateRepo(db *sql.DB, log *slog.Logger) repo.Template {
	return &templateRepo{
		db:  db,
		log: log,
	}
}

type repoTemplateScanner interface {
	Scan(dest ...any) error
}

func (r *templateRepo) FindByAlias(ctx context.Context, alias vo.Alias) (*entity.Template, error) {
	const query = `
		SELECT
			id,
			alias,
			subject,
			html,
			text
		FROM templates
		WHERE alias = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, alias.String())

	template, err := scanRepoTemplate(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		r.log.ErrorContext(ctx,
			"failed to find template by alias",
			slog.String("alias", alias.String()),
			slog.Any("error", err),
		)
		return nil, err
	}

	return template, nil
}

func (r *templateRepo) FindAll(ctx context.Context) ([]*entity.Template, error) {
	const query = `
		SELECT
			id,
			alias,
			subject,
			html,
			text
		FROM templates
		ORDER BY id
	`

	rows, err := executor(ctx, r.db).QueryContext(ctx, query)
	if err != nil {
		r.log.ErrorContext(ctx,
			"failed to find all templates",
			slog.Any("error", err),
		)
		return nil, err
	}
	defer rows.Close()

	templates := make([]*entity.Template, 0)

	for rows.Next() {
		template, err := scanRepoTemplate(rows)
		if err != nil {
			r.log.ErrorContext(ctx,
				"failed to scan template",
				slog.Any("error", err),
			)
			return nil, err
		}

		templates = append(templates, template)
	}

	if err := rows.Err(); err != nil {
		r.log.ErrorContext(ctx,
			"failed to iterate templates",
			slog.Any("error", err),
		)
		return nil, err
	}

	return templates, nil
}

func (r *templateRepo) Save(ctx context.Context, template *entity.Template) error {
	if template.IsDeleted() {
		const query = `
			DELETE FROM templates
			WHERE id = ?
		`

		if _, err := executor(ctx, r.db).ExecContext(ctx, query,
			template.ID().String(),
		); err != nil {
			r.log.ErrorContext(ctx,
				"failed to delete template",
				slog.String("template_id", template.ID().String()),
				slog.Any("error", err),
			)
			return err
		}
	} else if template.IsNew() {
		const query = `
			INSERT INTO templates (
				id,
				alias,
				subject,
				html,
				text
			)
			VALUES (?, ?, ?, ?, ?)
		`

		if _, err := executor(ctx, r.db).ExecContext(ctx, query,
			template.ID().String(),
			template.Alias().String(),
			template.Subject().String(),
			template.HTML().String(),
			template.Text().String(),
		); err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				switch {
				case strings.Contains(mysqlErr.Message, "uq_templates_alias"):
					return repo.ErrAliasAlreadyExists
				}
			}

			r.log.ErrorContext(ctx,
				"failed to insert template",
				slog.String("template_id", template.ID().String()),
				slog.Any("error", err),
			)
			return err
		}
	} else {
		const query = `
		    UPDATE templates
		    SET
		        alias = ?,
		        subject = ?,
		        html = ?,
		        text = ?
		    WHERE id = ?
		`

		if _, err := executor(ctx, r.db).ExecContext(ctx, query,
			template.Alias().String(),
			template.Subject().String(),
			template.HTML().String(),
			template.Text().String(),
			template.ID().String(),
		); err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				switch {
				case strings.Contains(mysqlErr.Message, "uq_templates_alias"):
					return repo.ErrAliasAlreadyExists
				}
			}

			r.log.ErrorContext(ctx,
				"failed to update template",
				slog.String("template_id", template.ID().String()),
				slog.Any("error", err),
			)
			return err
		}
	}

	template.MarkPersisted()

	return nil
}

func scanRepoTemplate(scanner repoTemplateScanner) (*entity.Template, error) {
	var (
		idRaw      string
		aliasRaw   string
		subjectRaw string
		htmlRaw    string
		textRaw    string
	)

	if err := scanner.Scan(
		&idRaw,
		&aliasRaw,
		&subjectRaw,
		&htmlRaw,
		&textRaw,
	); err != nil {
		return nil, err
	}

	id, _ := vo.NewTemplateID(idRaw)
	alias, _ := vo.NewAlias(aliasRaw)
	subject := vo.NewSubject(subjectRaw)
	html := vo.NewHTML(htmlRaw)
	text := vo.NewText(textRaw)

	return entity.ReconstituteTemplate(
		id,
		alias,
		subject,
		html,
		text,
	), nil
}
