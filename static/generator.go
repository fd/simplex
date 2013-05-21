package static

import (
	"database/sql"
	"simplex.sh/store"
)

type Generator interface {
	Generate(tx *Tx)
}

func Generate(src, dst store.Store, database *sql.DB, g Generator) error {
	tx := &Tx{
		src: src,
		dst: dst,
	}

	{ // database transaction
		db_tx, err := database.Begin()
		if err != nil {
			return err
		}

		tx.transaction = db_tx
		tx.database = database
	}

	g.Generate(tx)

	for _, t := range tx.terminators {
		tx.err.Add(t.Commit())
	}

	for _, t := range tx.terminators {
		tx.err.Add(t.Wait())
	}

	err := tx.err.Normalize()
	if err != nil {
		tx.transaction.Rollback()
		return err
	}

	return tx.transaction.Commit()
}

type GeneratorFunc func(tx *Tx)

func (f GeneratorFunc) Generate(tx *Tx) {
	f(tx)
}
