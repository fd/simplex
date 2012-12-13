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

	ctx.GolangFindFunctions()
	ctx.ParseTemplates()
	ctx.LookupFunctionCalls()
	ctx.UnfoldRenderFunctions()
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
	n = "\"github.com/fd/w/template/compiler/_test/app\".index"
	if tmpl, p := ctx.RenderFuncs[n]; !p {
		t.Errorf("missing template: %s", n)

		if tmpl.FunctionName() != "Index" {
			t.Errorf("expected Index but was %s", tmpl.FunctionName())
		}
	}

	n = "\"github.com/fd/w/template/compiler/_test/app\".index_1"
	if tmpl, p := ctx.RenderFuncs[n]; !p {
		t.Errorf("missing template: %s", n)

		if tmpl.FunctionName() != "index_1" {
			t.Errorf("expected index_1 but was %s", tmpl.FunctionName())
		}
	}

	n = "\"github.com/fd/w/template/compiler/_test/app\".helpers"
	if _, p := ctx.RenderFuncs[n]; p {
		t.Errorf("helpers.go is not a template: %s", n)
	}
}
