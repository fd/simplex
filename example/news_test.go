package example

import (
	"simplex.sh/static"
	"simplex.sh/stesting"
	"testing"
)

func TestSite(t *testing.T) {
	stesting.Golden(
		t,
		static.GeneratorFunc(Generate),
	)
}
