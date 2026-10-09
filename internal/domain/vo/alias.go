package vo

import "errors"

var (
	ErrAliasEmpty            = errors.New("alias must not be empty")
	ErrAliasTooLong          = errors.New("alias must not exceed 64 characters")
	ErrAliasInvalidCharacter = errors.New("alias contains an invalid character")
)

type Alias struct {
	value string
}

func NewAlias(value string) (Alias, error) {
	if value == "" {
		return Alias{}, ErrAliasEmpty
	}

	for i := 0; i < len(value); i++ {
		c := value[i]

		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '-' ||
			c == '_' ||
			c == '.') {
			return Alias{}, ErrAliasInvalidCharacter
		}
	}

	if len(value) > 64 {
		return Alias{}, ErrAliasTooLong
	}

	return Alias{
		value: value,
	}, nil
}

func (alias Alias) String() string {
	return alias.value
}

func (alias Alias) Equal(other Alias) bool {
	return alias.value == other.value
}
