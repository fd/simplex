package testing

import (
	"github.com/fd/simplex/data/storage/driver"
	T "testing"
)

func ValidateDriver(t *T.T, s driver.I) {
	var (
		dat  []byte
		err  error
		zero = [20]byte{}
	)

	// Get non existing object should return a driver.NotFound error
	dat, err = s.Get(zero)
	if err != driver.NotFound {
		t.Errorf("expected error to be driver.NotFound but was %v", err)
	}
	if dat != nil {
		t.Errorf("expected dat to be nil but was %+v", dat)
	}

	// set a value
	err = s.Set(zero, []byte("foo"))
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}

	// get a existing value
	dat, err = s.Get(zero)
	if err != nil {
		t.Errorf("expected error to be nil but was %v", err)
	}
	if string(dat) != "foo" {
		t.Errorf("expected data to be 'foo' but was %+v", dat)
	}
}
