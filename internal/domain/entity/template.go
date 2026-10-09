package entity

import (
	"errors"

	"github.com/velonyapp/email/internal/domain/vo"
)

var (
	ErrTemplateBodyRequired = errors.New("at least one template body is required")
)

type Template struct {
	id      vo.TemplateID
	alias   vo.Alias
	subject vo.Subject
	html    vo.HTML
	text    vo.Text
	deleted bool
}

func NewTemplate(
	alias vo.Alias,
	subject vo.Subject,
	html vo.HTML,
	text vo.Text,
) (*Template, error) {
	if html.IsEmpty() && text.IsEmpty() {
		return nil, ErrTemplateBodyRequired
	}

	return &Template{
		id:      vo.GenerateTemplateID(),
		alias:   alias,
		subject: subject,
		html:    html,
		text:    text,
	}, nil
}

func ReconstituteTemplate(
	id vo.TemplateID,
	alias vo.Alias,
	subject vo.Subject,
	html vo.HTML,
	text vo.Text,
) (*Template, error) {
	t := &Template{
		id:      id,
		alias:   alias,
		subject: subject,
		html:    html,
		text:    text,
	}

	if !t.HasHTML() && !t.HasText() {
		return nil, ErrTemplateBodyRequired
	}

	return t, nil
}

func (t *Template) ID() vo.TemplateID {
	return t.id
}

func (t *Template) Alias() vo.Alias {
	return t.alias
}

func (t *Template) Subject() vo.Subject {
	return t.subject
}

func (t *Template) HTML() vo.HTML {
	return t.html
}

func (t *Template) Text() vo.Text {
	return t.text
}

func (t *Template) IsDeleted() bool {
	return t.deleted
}

func (t *Template) HasHTML() bool {
	return !t.html.IsEmpty()
}

func (t *Template) HasText() bool {
	return !t.text.IsEmpty()
}

func (t *Template) ChangeAlias(newAlias vo.Alias) {
	t.alias = newAlias
}

func (t *Template) ChangeSubject(newSubject vo.Subject) {
	t.subject = newSubject
}

func (t *Template) ChangeHTML(newHTML vo.HTML) error {
	if newHTML.IsEmpty() && !t.HasText() {
		return ErrTemplateBodyRequired
	}

	t.html = newHTML

	return nil
}

func (t *Template) ChangeText(newText vo.Text) error {
	if newText.IsEmpty() && !t.HasHTML() {
		return ErrTemplateBodyRequired
	}

	t.text = newText

	return nil
}

func (t *Template) ChangeHTMLAndText(newHTML vo.HTML, newText vo.Text) error {
	if newHTML.IsEmpty() && newText.IsEmpty() {
		return ErrTemplateBodyRequired
	}

	t.html = newHTML
	t.text = newText

	return nil
}

func (t *Template) Delete() {
	t.deleted = true
}
