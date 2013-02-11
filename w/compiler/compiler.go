package compiler

import (
	container "simplex.sh/w/container/compiler"
	template "simplex.sh/w/template/compiler"
)

func Compile(pwd string) error {
	var err error

	ctx := template.NewContext(pwd)

	// compile applications
	err = ctx.Compile()
	if err != nil {
		return err
	}

	// compile schemas           (to go files)
	// compile templates         (to go files)
	err = ctx.WriteTemplateFiles()
	if err != nil {
		return err
	}
	defer ctx.RemoveTemplateFiles()

	// compile container wrapper (to go files)
	err = container.WriteContainerWrapper(ctx)
	if err != nil {
		return err
	}
	defer container.RemoveContainerWrapper(ctx)

	// compile additional tests  (to go files)
	// compile container         (to binary)
	err = compile_go_files(ctx)
	if err != nil {
		return err
	}

	// run tests

	/*
	   fmt.Printf("Helpers:\n")
	   for n := range ctx.Helpers {
	     fmt.Printf("- %s\n", n)
	   }

	   fmt.Printf("DataViews:\n")
	   for n := range ctx.DataViews {
	     fmt.Printf("- %s\n", n)
	   }

	   fmt.Printf("RenderFuncs:\n")
	   for n := range ctx.RenderFuncs {
	     fmt.Printf("- %s\n", n)
	   }
	*/

	return nil
}
