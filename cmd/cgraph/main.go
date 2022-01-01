package main

import (
	"cgraph"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(cgraph.Analyzer) }
