package graph

import (
	"fmt"
	"io"
	"strconv"
)

// DotNode is an interface to describe the DOT requirements for traversing graphs
type DotNode interface {
	Label() string
	GetGID() int
	GetGCost() int64
}

// DotWriter is a type that will build the DOT syntax of our Metadata
type DotWriter struct {
	output   io.Writer
	MaxDepth int
	Reversed bool
}

type plotCtx struct {
	nodeFlags map[string]bool
	edgeFlags map[string]bool
	level     int
}

func (ctx *plotCtx) isPlottedNode(node Relationship) bool {
	_, ok := ctx.nodeFlags[node.GetID()]
	return ok
}
func (ctx *plotCtx) setPlotted(node Relationship) {
	_, ok := ctx.nodeFlags[node.GetID()]
	if !ok {
		ctx.nodeFlags[node.GetID()] = true
	}
}

func (ctx *plotCtx) isDepthOver() bool {
	return (ctx.level <= 0)
}

func (ctx *plotCtx) Deeper() *plotCtx {
	return &plotCtx{
		nodeFlags: ctx.nodeFlags,
		edgeFlags: ctx.edgeFlags,
		level:     ctx.level - 1,
	}
}

func newPlotContext(level int) *plotCtx {
	return &plotCtx{
		level:     level,
		nodeFlags: make(map[string]bool),
		edgeFlags: make(map[string]bool),
	}
}

func (ctx *plotCtx) isPlottedEdge(nodeA, nodeB Relationship) bool {
	edgeName := fmt.Sprintf("%s->%s", nodeA.GetID(), nodeB.GetID())
	_, ok := ctx.edgeFlags[edgeName]
	if !ok {
		ctx.edgeFlags[edgeName] = true
	}
	return ok
}

// NewDotWriter creates a new writer to generate DOT graphs
func NewDotWriter(output io.Writer, maxdepth int, reversed bool) *DotWriter {
	return &DotWriter{
		output:   output,
		MaxDepth: maxdepth,
		Reversed: reversed,
	}
}

// PlotGraph traverses the root Relationship (root) while DotWriter.MaxDepth > 0, generating the graph.
func (dw *DotWriter) PlotGraph(root Relationship) {
	dw.printLine("digraph main{")
	dw.printLine("\tedge[arrowhead=vee]")
	dw.printLine("\tgraph [rankdir=LR,compound=true,ranksep=1.0];")
	dw.plotNode(newPlotContext(dw.MaxDepth), root)
	dw.printLine("}")
}

func (dw *DotWriter) plotNode(ctx *plotCtx, node Relationship) {
	if ctx.isPlottedNode(node) {
		return
	}
	if ctx.isDepthOver() {
		return
	}
	ctx.setPlotted(node)
	dw.plotNodeStyle(node)
	for _, s := range dw.getDependency(node) {
		dw.plotEdge(ctx, node, s)
		dw.plotNode(ctx.Deeper(), s)
	}
}

func (dw *DotWriter) getDependency(node Relationship) []Relationship {
	if dw.Reversed {
		return node.Parents()
	}
	return node.Children()
}

func (dw *DotWriter) plotNodeStyle(node Relationship) {
	dw.printFormat("\t/* plot %s */\n", node.GetID())
	dw.printFormat("\t%s[shape=%s,label=\"%s\",style=%s]\n",
		escape(node.GetID()),
		escape("record"),
		node.Label(),
		escape("solid"),
	)
}

func (dw *DotWriter) plotEdge(ctx *plotCtx, nodeA, nodeB Relationship) {
	if ctx.isPlottedEdge(nodeA, nodeB) {
		return
	}
	dir := "forward"
	if dw.Reversed {
		dir = "back"
	}
	dw.printFormat("\t%s -> %s[dir=%s]\n", escape(nodeA.GetID()), escape(nodeB.GetID()), dir)
}

func (dw *DotWriter) printLine(str string) {
	fmt.Fprintln(dw.output, str)
}

func (dw *DotWriter) printFormat(pattern string, args ...interface{}) {
	fmt.Fprintf(dw.output, pattern, args...)
}

func escape(target string) string {
	return strconv.Quote(target)
}
