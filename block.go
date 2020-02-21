package flub

import (
	"fmt"

	"github.com/influxdata/flux/ast"
)

type Block struct {
	nodes    []*node
	vars     map[string]ast.Expression
	varnames []string
	id       int
}

func newBlock() *Block {
	return &Block{
		vars: make(map[string]ast.Expression),
	}
}

func (b *Block) Call(name string, args ...*KeyValuePair) Pipeline {
	expr := Call(name, args...)
	return b.newPipeline(expr)
}

func (b *Block) Eval(n Node) Pipeline {
	return Pipeline{n: b.newNode(func() ast.Expression {
		return n.get()
	})}
}

func (b *Block) From(bucket string) Pipeline {
	return b.Call("from",
		KV("bucket", String(bucket)),
	)
}

func (b *Block) newNode(create func() ast.Expression) *node {
	n := &node{b: b, create: create}
	b.nodes = append(b.nodes, n)
	return n
}

func (b *Block) newPipeline(start ast.Expression) Pipeline {
	return Pipeline{n: b.newNode(func() ast.Expression {
		return start
	})}
}

func (b *Block) build() []ast.Statement {
	// Instantiate the body backwards.
	var body []ast.Statement
	for i := len(b.nodes) - 1; i >= 0; i-- {
		n := b.nodes[i]
		if n.children == 0 {
			body = append(body, &ast.ExpressionStatement{
				Expression: n.get(),
			})
		} else if n.children > 1 || n.varname != "" {
			body = append(body, &ast.VariableAssignment{
				ID:   &ast.Identifier{Name: n.varname},
				Init: b.vars[n.varname],
			})
		}
	}

	for i, j := 0, len(body)-1; i < j; i, j = i+1, j-1 {
		body[i], body[j] = body[j], body[i]
	}
	return body
}

func (b *Block) genVarName() string {
	for {
		id := fmt.Sprintf("var%d", b.id)
		b.id++
		if _, ok := b.vars[id]; !ok {
			return id
		}
	}
}
