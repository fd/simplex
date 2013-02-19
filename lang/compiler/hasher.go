package compiler

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"simplex.sh/lang/types"
	"sort"
)

func (ctx *Context) generate_version_hash() error {
	sha := sha1.New()

	for _, name := range ctx.GoFiles {
		data, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		io.WriteString(sha, name)
		sha.Write(data)
	}

	for _, name := range ctx.SxFiles {
		data, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		io.WriteString(sha, name)
		sha.Write(data)
	}

	ids := make([]string, 0, len(ctx.TypesPackage.Imports))
	for id, _ := range ctx.TypesPackage.Imports {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		pkg := ctx.TypesPackage.Imports[id]
		obj := pkg.Scope.Lookup("SxVersion")
		if vsn, ok := obj.(*types.Const); ok && vsn != nil {
			if vsn_str, ok := vsn.Val.(string); ok && vsn_str != "" {
				io.WriteString(sha, id)
				io.WriteString(sha, vsn_str)
				continue
			}
		}

		filename, id := types.FindPkg(id, ctx.TypesPackage.Path)
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		io.WriteString(sha, id)
		sha.Write(data)
	}

	ctx.Version = hex.EncodeToString(sha.Sum(nil))
	return nil
}
