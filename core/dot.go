package core

import (
	"fmt"
	"io"
	"strconv"
)

// DotNode is an interface to describe the DOT requirements for traversing graphs
type DotNode interface {
	Name() string
	Label() string
	Shape() string
	Style() string
	GetID() string
	GetGID() int
	GetGCost() int64
	GetChecksum() uint64
	Children() []DotNode
	Parents() []DotNode
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

func (ctx *plotCtx) isPlottedNode(node DotNode) bool {
	_, ok := ctx.nodeFlags[node.Name()]
	return ok
}
func (ctx *plotCtx) setPlotted(node DotNode) {
	_, ok := ctx.nodeFlags[node.Name()]
	if !ok {
		ctx.nodeFlags[node.Name()] = true
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

func (ctx *plotCtx) isPlottedEdge(nodeA, nodeB DotNode) bool {
	edgeName := fmt.Sprintf("%s->%s", nodeA.Name(), nodeB.Name())
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

// PlotGraph traverses the root DotNode (root) while DotWriter.MaxDepth > 0, generating the graph.
func (dw *DotWriter) PlotGraph(root DotNode) {
	dw.printLine("digraph main{")
	dw.printLine("\tedge[arrowhead=vee]")
	dw.printLine("\tgraph [rankdir=LR,compound=true,ranksep=1.0];")
	dw.plotNode(newPlotContext(dw.MaxDepth), root)
	dw.printLine("}")
}

func (dw *DotWriter) plotNode(ctx *plotCtx, node DotNode) {
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

func (dw *DotWriter) getDependency(node DotNode) []DotNode {
	if dw.Reversed {
		return node.Parents()
	}
	return node.Children()
}

func (dw *DotWriter) plotNodeStyle(node DotNode) {
	dw.printFormat("\t/* plot %s */\n", node.Name())
	dw.printFormat("\t%s[shape=%s,label=\"%s\",style=%s]\n",
		escape(node.Name()),
		escape(node.Shape()),
		node.Label(),
		escape(node.Style()),
	)
}

func (dw *DotWriter) plotEdge(ctx *plotCtx, nodeA, nodeB DotNode) {
	if ctx.isPlottedEdge(nodeA, nodeB) {
		return
	}
	dir := "forward"
	if dw.Reversed {
		dir = "back"
	}
	dw.printFormat("\t%s -> %s[dir=%s]\n", escape(nodeA.Name()), escape(nodeB.Name()), dir)
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
