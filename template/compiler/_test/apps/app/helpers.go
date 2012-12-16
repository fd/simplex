package app

import (
	"github.com/fd/w/data"
	"strings"
)

func LinkTo(title, url string) data.Value {
	return strings.ToTitle("(LINK)")
}

func Tag(name string, attrs map[string]string) string {
	return "(TAG)"
}
func Tag1(name string, attrs map[string]string) (data.Value, error) {
	return "(TAG)", nil
}

func Tag2(name string, attrs map[string]string) data.Value {
	return "(TAG)"
}
