package vo

type Text struct {
	value string
}

func NewText(value string) Text {
	return Text{value: value}
}

func (t Text) String() string {
	return t.value
}

func (t Text) Equal(other Text) bool {
	return t.value == other.value
}

func (t Text) IsEmpty() bool {
	return t.value == ""
}
