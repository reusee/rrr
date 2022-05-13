package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {

	// load
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedDeps |
			packages.NeedFiles |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedName |
			packages.NeedSyntax,
		Tests: true,
	}, os.Args[1:]...)
	ce(err)
	if packages.PrintErrors(pkgs) > 0 {
		return
	}

	// rename identifiers
	renameIdentifier(pkgs, "clusterID", "shardID")
	renameIdentifier(pkgs, "ClusterID", "ShardID")

	// save all files
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			position := pkg.Fset.Position(file.Pos())
			if !strings.HasSuffix(position.Filename, ".go") {
				continue
			}
			buf := new(bytes.Buffer)
			ce(format.Node(buf, pkg.Fset, file))
			ce(os.WriteFile(position.Filename, buf.Bytes(), 0644))
		}
	}

}

func renameIdentifier(pkgs []*packages.Package, oldName string, newName string) {
	identifiers := make(map[*ast.Ident]bool)

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {

			ast.Inspect(file, func(node ast.Node) bool {
				ident, ok := node.(*ast.Ident)
				if !ok {
					return true
				}

				if ident.Name != oldName {
					return true
				}

				// defined in current package
				obj, ok := pkg.TypesInfo.Defs[ident]
				if !ok || obj == nil {
					return true
				}
				identifiers[ident] = true

				// iter all packages
				for _, pkg := range pkgs {
					for id, o := range pkg.TypesInfo.Uses {
						if o != obj {
							continue
						}
						identifiers[id] = true
					}
				}

				return true
			})

		}
	}

	for id := range identifiers {
		id.Name = newName
	}
}
