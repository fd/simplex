package compiler

import (
	"os"
	"testing"
)

func TestContext(t *testing.T) {

	pwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	ctx := &Context{
		WROOT: pwd,
	}

	err = ctx.Analyze("./_test/app")
	if err != nil {
		t.Error(err)
	}

	ctx.CleanTemplates()

	var n string

	// check the helpers
	n = "\"github.com/fd/w/template/compiler/_test/app\".LinkTo"
	if _, p := ctx.Helpers[n]; !p {
		t.Errorf("missing helper: %s", n)
	}

	n = "\"github.com/fd/w/template/compiler/_test/app\".Tag"
	if _, p := ctx.Helpers[n]; !p {
		t.Errorf("missing helper: %s", n)
	}

	n = "\"github.com/fd/w/template/compiler/_test/app\".AllPosts"
	if _, p := ctx.Helpers[n]; p {
		t.Errorf("data.View is not a helper: %s", n)
	}

	// check the templates
	n = "\"github.com/fd/w/template/compiler/_test/app\".Template(index.go.html:0)"
	if _, p := ctx.Templates[n]; !p {
		t.Errorf("missing template: %s", n)
	}

	n = "\"github.com/fd/w/template/compiler/_test/app\".Template(helpers.go:0)"
	if _, p := ctx.Templates[n]; p {
		t.Errorf("helpers.go is not a template: %s", n)
	}
}
