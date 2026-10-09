package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/velonyapp/email/internal/domain/entity"
	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"
)

var _ repo.Template = (*templateRepo)(nil)

type templateRepo struct {
	db *sql.DB
}

func NewTemplateRepo(
	db *sql.DB,
) repo.Template {
	return &templateRepo{
		db: db,
	}
}

type templateScanner interface {
	Scan(dest ...any) error
}

func (r *templateRepo) FindByID(ctx context.Context, templateID vo.TemplateID) (*entity.Template, error) {
	const query = `
		SELECT
			id,
			subject,
			html,
			text
		FROM templates
		WHERE id = ?
		LIMIT 1
		FOR UPDATE
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query,
		templateID.String(),
	)

	template, err := scanTemplate(row)
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
			subject,
			html,
			text
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
		template, err := scanTemplate(rows)
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
	} else {
		const query = `
			INSERT INTO templates (
				id,
				subject,
				html,
				text
			)
			VALUES (?, ?, ?, ?) AS new
			ON DUPLICATE KEY UPDATE
				subject = new.subject,
				html = new.html,
				text = new.text
		`

		if _, err := executor(ctx, r.db).ExecContext(ctx, query,
			template.ID().String(),
			template.Subject().String(),
			template.HTML().String(),
			template.Text().String(),
		); err != nil {
			return err
		}
	}

	return nil
}

func scanTemplate(scanner templateScanner) (*entity.Template, error) {
	var (
		idRaw      string
		subjectRaw string
		htmlRaw    string
		textRaw    string
	)

	if err := scanner.Scan(
		&idRaw,
		&subjectRaw,
		&htmlRaw,
		&textRaw,
	); err != nil {
		return nil, err
	}

	id, err := vo.NewTemplateID(idRaw)
	if err != nil {
		return nil, err
	}
	subject := vo.NewSubject(subjectRaw)
	html := vo.NewHTML(htmlRaw)
	text := vo.NewText(textRaw)

	return entity.ReconstituteTemplate(
		id,
		subject,
		html,
		text,
	)
}
