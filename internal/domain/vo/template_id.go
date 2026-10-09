package vo

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTemplateIDInvalid = errors.New("invalid template ID")
)

type TemplateID struct {
	value string
}

func NewTemplateID(value string) (TemplateID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return TemplateID{}, ErrTemplateIDInvalid
	}

	return TemplateID{value: id.String()}, nil
}

func GenerateTemplateID() TemplateID {
	return TemplateID{
		value: uuid.Must(uuid.NewV7()).String(),
	}
}

func (id TemplateID) String() string {
	return id.value
}

func (id TemplateID) Equal(other TemplateID) bool {
	return id.value == other.value
}
