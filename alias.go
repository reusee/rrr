package main

import (
	"fmt"

	"github.com/reusee/e4"
)

var (
	ce = e4.Check.With(e4.WrapStacktrace)
	pt = fmt.Printf
)
