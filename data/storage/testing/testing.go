package testing

import (
	"github.com/fd/w/data/storage/raw/driver"
	"strings"
	T "testing"
)

func ValidateRawDriver(t *T.T, s driver.I) {
	var (
		dat []byte
		ids []string
		err error
	)

	// ids should be empty
	ids, err = s.Ids()
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected len(ids) to be 0 but was %+v", len(ids))
	}

	// get a none existent value
	dat, err = s.Get("a")
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if dat != nil {
		t.Errorf("expected data to be nil but was %+v", dat)
	}

	// set a new value
	err = s.Commit(
		map[string][]byte{
			"a": []byte("foo"),
			"b": []byte("bar"),
		},
		nil,
	)
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}

	// ids should be [a b]
	ids, err = s.Ids()
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if strings.Join(ids, ",") != "a,b" {
		t.Errorf("expected ids to be [a b] but was %+v", ids)
	}

	// get a existing value
	dat, err = s.Get("a")
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if string(dat) != "foo" {
		t.Errorf("expected data to be 'foo' but was %+v", dat)
	}

	// get a existing value
	dat, err = s.Get("b")
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if string(dat) != "bar" {
		t.Errorf("expected data to be 'bar' but was %+v", dat)
	}

	// update a value
	err = s.Commit(
		map[string][]byte{
			"b": []byte("baz"),
			"c": []byte("bax"),
		},
		nil,
	)
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}

	// ids should be [a b c]
	ids, err = s.Ids()
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if strings.Join(ids, ",") != "a,b,c" {
		t.Errorf("expected ids to be [a b c] but was %+v", ids)
	}

	// get a existing value
	dat, err = s.Get("a")
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if string(dat) != "foo" {
		t.Errorf("expected data to be 'foo' but was %+v", dat)
	}

	// get a existing value
	dat, err = s.Get("b")
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if string(dat) != "baz" {
		t.Errorf("expected data to be 'bar' but was %+v", dat)
	}

	// get a existing value
	dat, err = s.Get("c")
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if string(dat) != "bax" {
		t.Errorf("expected data to be 'bar' but was %+v", dat)
	}

	// update a value
	err = s.Commit(
		map[string][]byte{},
		[]string{"a", "b"},
	)
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}

	// ids should be [c]
	ids, err = s.Ids()
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if strings.Join(ids, ",") != "c" {
		t.Errorf("expected ids to be [c] but was %+v", ids)
	}

	// get a existing value
	dat, err = s.Get("a")
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if dat != nil {
		t.Errorf("expected data to be nil but was %+v", dat)
	}

	// get a existing value
	dat, err = s.Get("b")
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if dat != nil {
		t.Errorf("expected data to be nil but was %+v", dat)
	}

	// get a existing value
	dat, err = s.Get("c")
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if string(dat) != "bax" {
		t.Errorf("expected data to be 'bar' but was %+v", dat)
	}

}
