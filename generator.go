package static

import (
	"github.com/fd/static/store"
)

type Generator interface {
	Generate(tx *Tx)
}

type GeneratorFunc func(tx *Tx)

func Generate(src, dst store.Store, g Generator) error {
	tx := &Tx{
		src: src,
		dst: dst,
	}

	g.Generate(tx)

	for _, t := range tx.terminators {
		err := t.Commit()
		if err != nil {
			tx.err.Add(err)
		}
	}

	for _, t := range tx.terminators {
		err := t.Wait()
		if err != nil {
			tx.err.Add(err)
		}
	}

	return tx.err.Normalize()
}

func (f GeneratorFunc) Generate(tx *Tx) {
	f(tx)
}
