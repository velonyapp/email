package vo

type HTML struct {
	value string
}

func NewHTML(value string) HTML {
	return HTML{value: value}
}

func (h HTML) String() string {
	return h.value
}

func (h HTML) Equal(other HTML) bool {
	return h.value == other.value
}

func (h HTML) IsEmpty() bool {
	return h.value == ""
}
