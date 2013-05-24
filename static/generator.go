package static

import (
	"database/sql"
	"simplex.sh/store"
	"simplex.sh/store/cas"
	"simplex.sh/store/router"
)

type Generator interface {
	Generate(tx *Tx)
}

func Generate(src store.Store, database *sql.DB, g Generator) error {
	tx := &Tx{
		src: src,
	}

	{ // database transaction
		db_tx, err := database.Begin()
		if err != nil {
			return err
		}

		tx.database = database
		tx.transaction = db_tx
	}

	{ // cas transaction
		cas_writer, err := cas.OpenWriter(tx.transaction)
		if err != nil {
			return err
		}

		tx.cas_writer = cas_writer
	}

	{
		w, err := router.OpenWriter(tx.transaction)
		if err != nil {
			return err
		}

		tx.router_writer = w
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
