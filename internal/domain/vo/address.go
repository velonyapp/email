package vo

import (
	"errors"
	"net/mail"
	"strings"
)

var (
	ErrAddressInvalid = errors.New("address is invalid")
)

type Address struct {
	value string
}

func NewAddress(value string) (Address, error) {
	value = strings.TrimSpace(value)

	if _, err := mail.ParseAddress(value); err != nil {
		return Address{}, ErrAddressInvalid
	}

	return Address{
		value: value,
	}, nil
}

func (s Address) String() string {
	return s.value
}

func (s Address) Equal(other Address) bool {
	return s.value == other.value
}
