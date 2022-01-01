package cgraph

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/callgraph/vta"
	"golang.org/x/tools/go/ssa/ssautil"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
)

const doc = "cgraph is ..."

type calleesFact map[types.Object]struct{}

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "cgraph",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		buildssa.Analyzer,
	},
	//FactTypes: []analysis.Fact{},
}

func GraphVisitEdges1(g *callgraph.Graph, edge func(*callgraph.Edge) error) error {

	seen := make(map[*callgraph.Node]bool)
	var visit func(n *callgraph.Node) error
	visit = func(n *callgraph.Node) error {
		fmt.Println(n.ID)

		if !seen[n] {
			seen[n] = true
			for _, e := range n.In {
				if err := visit(e.Callee); err != nil {
					return err
				}
				if err := edge(e); err != nil {
					return err
				}
			}
		}
		return nil
	}
	for _, n := range g.Nodes {
		if err := visit(n); err != nil {
			return err
		}
	}
	return nil
}

type UnionFind struct {
	par []int
}

func newUnionFind(n int) UnionFind {
	par := make([]int, n+1)
	for i := 0; i < len(par); i++ {
		par[i] = i
	}
	return UnionFind{par}
}

func (u *UnionFind) root(x int) int {
	if u.par[x] < 0 {
		return x
	}
	u.par[x] = u.root(u.par[x])
	return u.par[x]
}

func (u *UnionFind) unite(x, y int) bool {
	x, y = u.root(x), u.root(y)
	if x == y {
		return false
	}
	u.par[x] = y
	return true
}

func run(pass *analysis.Pass) (interface{}, error) {

	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	initcg := cha.CallGraph(s.Pkg.Prog)
	cg := vta.CallGraph(ssautil.AllFunctions(s.Pkg.Prog), initcg)

	var g map[string][]string
	g = make(map[string][]string)
	var edges []string
	callgraph.GraphVisitEdges(cg, func(edge *callgraph.Edge) error {
		caller := edge.Caller.Func.String()
		callee := edge.Callee.Func.String()
		if _, ok := g[callee]; !ok {
			g[callee] = make([]string, 0)
		}
		g[callee] = append(g[callee], caller)

		return nil
	})

	seen := make(map[string]bool)
	var visit func(string) error
	visit = func(fn string) error {
		fmt.Println(fn)
		if !seen[fn] {
			seen[fn] = true
			for _, e := range g[fn] {
				visit(e)
			}
		}
		return nil
	}
	visit("a.g")

	return nil, nil
}
