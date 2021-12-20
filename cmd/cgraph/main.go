package main

import (
	"cgraph"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(cgraph.Analyzer) }
