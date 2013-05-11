package example

import (
	"github.com/fd/static"
	"github.com/fd/static/stesting"
	"testing"
)

func TestSite(t *testing.T) {
	stesting.Golden(
		t,
		static.GeneratorFunc(Generate),
	)
}
