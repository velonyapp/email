package vo

import "errors"

var (
	ErrTemplateIDEmpty            = errors.New("template ID must not be empty")
	ErrTemplateIDTooLong          = errors.New("template ID must not exceed 64 characters")
	ErrTemplateIDInvalidCharacter = errors.New("template ID contains an invalid character")
)

type TemplateID struct {
	value string
}

func NewTemplateID(value string) (TemplateID, error) {
	if value == "" {
		return TemplateID{}, ErrTemplateIDEmpty
	}

	for i := 0; i < len(value); i++ {
		c := value[i]

		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '-' ||
			c == '_' ||
			c == '.') {
			return TemplateID{}, ErrTemplateIDInvalidCharacter
		}
	}

	if len(value) > 64 {
		return TemplateID{}, ErrTemplateIDTooLong
	}

	return TemplateID{
		value: value,
	}, nil
}

func (id TemplateID) String() string {
	return id.value
}

func (id TemplateID) Equal(other TemplateID) bool {
	return id.value == other.value
}
