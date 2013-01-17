package runtime

import (
	"fmt"
)

type HTML string

type Buffer struct {
	instructions []Instruction
}

type Instruction interface {
	ToHtmlString() string
}

func (b *Buffer) Write(v interface{}) {
	switch l := v.(type) {

	case string:
		// escape HTML
		b.instructions = append(b.instructions, HTML(l))

	case HTML:
		b.instructions = append(b.instructions, l)

	default:
		panic(fmt.Sprintf("%T doesn't implement the runtime.Instruction interface.", v))

	}
}

func (h HTML) ToHtmlString() string {
	return string(h)
}
