package collate

type (
	String string
	Bytea  []byte
	Runea  []rune

	localized interface {
		is_localized()
	}
)

func (String) is_localized() {}
func (Bytea) is_localized()  {}
func (Runea) is_localized()  {}
