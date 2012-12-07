package compiler

import (
	"os"
	"testing"
)

func TestContect(t *testing.T) {

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = Context(pwd + "/../../apps/cp")
	if err != nil {
		t.Fatal(err)
	}

}
