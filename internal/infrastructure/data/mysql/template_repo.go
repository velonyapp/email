package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/velonyapp/email/internal/domain/entity"
	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"

	"github.com/go-sql-driver/mysql"
)

var _ repo.Template = (*templateRepo)(nil)

type templateRepo struct {
	db *sql.DB
}

func NewTemplateRepo(db *sql.DB) repo.Template {
	return &templateRepo{
		db: db,
	}
}

type repoTemplateScanner interface {
	Scan(dest ...any) error
}

func (r *templateRepo) FindByID(ctx context.Context, templateID vo.TemplateID) (*entity.Template, error) {
	const query = `
		SELECT
			id,
			alias,
			subject,
			html,
			text,
			version
		FROM templates
		WHERE id = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, templateID.String())

	template, err := scanRepoTemplate(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return template, nil
}

func (r *templateRepo) FindByAlias(ctx context.Context, alias vo.Alias) (*entity.Template, error) {
	const query = `
		SELECT
			id,
			alias,
			subject,
			html,
			text,
			version
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
			text,
			version
		FROM templates
		ORDER BY id
	`

	rows, err := executor(ctx, r.db).QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := make([]*entity.Template, 0)

	for rows.Next() {
		template, err := scanRepoTemplate(rows)
		if err != nil {
			return nil, err
		}

		templates = append(templates, template)
	}

	if err := rows.Err(); err != nil {
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
			return err
		}
	} else if template.IsNew() {
		const query = `
			INSERT INTO templates (
				id,
				alias,
				subject,
				html,
				text,
				version
			)
			VALUES (?, ?, ?, ?, ?, ?)
		`

		if _, err := executor(ctx, r.db).ExecContext(ctx, query,
			template.ID().String(),
			template.Alias().String(),
			template.Subject().String(),
			template.HTML().String(),
			template.Text().String(),
			template.Version()+1,
		); err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				switch {
				case strings.Contains(mysqlErr.Message, "uq_templates_alias"):
					return repo.ErrAliasAlreadyExists
				}
			}

			return err
		}
	} else {
		const query = `
		    UPDATE templates
		    SET
		        alias = ?,
		        subject = ?,
		        html = ?,
		        text = ?,
		        version = version + 1
		    WHERE id = ?
		        AND version = ?
		`

		result, err := executor(ctx, r.db).ExecContext(ctx, query,
			template.Alias().String(),
			template.Subject().String(),
			template.HTML().String(),
			template.Text().String(),
			template.ID().String(),
			template.Version(),
		)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				switch {
				case strings.Contains(mysqlErr.Message, "uq_templates_alias"):
					return repo.ErrAliasAlreadyExists
				}
			}

			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return repo.ErrTemplateVersionConflict
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
		versionRaw uint32
	)

	if err := scanner.Scan(
		&idRaw,
		&aliasRaw,
		&subjectRaw,
		&htmlRaw,
		&textRaw,
		&versionRaw,
	); err != nil {
		return nil, err
	}

	id, _ := vo.NewTemplateID(idRaw)
	alias, _ := vo.NewAlias(aliasRaw)
	subject := vo.NewSubject(subjectRaw)
	html := vo.NewHTML(htmlRaw)
	text := vo.NewText(textRaw)
	version := versionRaw

	return entity.ReconstituteTemplate(
		id,
		alias,
		subject,
		html,
		text,
		version,
	), nil
}
