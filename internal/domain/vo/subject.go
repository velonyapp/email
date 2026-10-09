package vo

import "strings"

type Subject struct {
	value string
}

func NewSubject(value string) Subject {
	return Subject{
		value: strings.TrimSpace(value),
	}
}

func (s Subject) String() string {
	return s.value
}

func (s Subject) Equal(other Subject) bool {
	return s.value == other.value
}
