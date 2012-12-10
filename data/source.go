package data

import (
	"github.com/fd/w/data/storage"
)

type Source struct {
	Storage storage.I
}

func (s *Source) Get(id string) Value {
	var (
		dat []byte
		raw interface{}
		val Value
		err error
	)

	dat, err = s.Storage.Get(id)
	if err != nil {
		// log err
		return nil
	}

	err = json.Unmarshal(data, &raw)
	if err != nil {
		// log err
		return nil
	}

	val, err = json_to_value(raw)
	if err != nil {
		// log err
		return nil
	}

	return val
}

func (s *Source) Ids() []string {
	var (
		ids []string
		err error
	)

	ids, err = s.Storage.Ids()
	if err != nil {
		// log err
		return nil
	}

	return ids
}

func (s *Source) Commit(set map[string]Value, del []string) {

}
