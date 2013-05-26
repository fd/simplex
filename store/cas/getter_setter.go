package cas

import (
	"database/sql"
)

type GetterSetter interface {
	Getter
	Setter
}

type getter_setter_t struct {
	getter_t
	setter_t
}

func Open(db *sql.DB) (GetterSetter, error) {
	gs := &getter_setter_t{}

	err := gs.init(db)
	if err != nil {
		return nil, err
	}

	err = update_schema(db)
	if err != nil {
		return nil, err
	}

	return gs, nil
}

func (gs *getter_setter_t) init(db *sql.DB) error {
	err := gs.getter_t.init(db)
	if err != nil {
		return err
	}

	err = gs.setter_t.init(db)
	if err != nil {
		return err
	}

	return nil
}
