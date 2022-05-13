package main

import (
	"os"

	"golang.org/x/tools/go/packages"
)

func main() {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedDeps |
			packages.NeedFiles |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedSyntax,
		Tests: true,
	}, os.Args[1:]...)
	ce(err)
	if packages.PrintErrors(pkgs) {
		return
	}
}
