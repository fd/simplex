package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

type List []error

func (l List) Normalize() error {
	if l.HasErrors() {
		return l
	}

	return nil
}

func (l List) HasErrors() bool {
	return len(l) > 0
}

func (l *List) Add(err error) {
	if err == nil {
		return
	}

	if list, ok := err.(List); ok {
		for _, err := range list {
			if err != nil {
				*l = append(*l, err)
			}
		}
		return
	}

	*l = append(*l, err)
}

func (l *List) Errorf(format string, args ...interface{}) {
	l.Add(fmt.Errorf(format, args...))
}

func (l List) Error() string {
	s := make([]string, len(l))

	for i, err := range l {
		s[i] = err.Error()
	}

	return strings.Join(s, "\n")
}

func (l List) MarshalJSON() ([]byte, error) {
	if !l.HasErrors() {
		return []byte("null"), nil
	}

	s := make([]interface{}, len(l))

	for i, err := range l {
		if _, ok := err.(json.Marshaler); ok {
			s[i] = err
		} else {
			s[i] = err.Error()
		}
	}

	return json.Marshal(s)
}
