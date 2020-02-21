package flub

import "github.com/influxdata/flux/ast"

type Node interface {
	node(b *Block) *node
	get() ast.Expression
}

type node struct {
	b        *Block
	expr     ast.Expression
	create   func() ast.Expression
	varname  string
	children int
}

func (n *node) node(b *Block) *node {
	if b == nil {
		return n
	} else if n.b == nil {
		// This node is not associated with a block
		// so duplicate the contents and assign
		// the block to the duplicate.
		nn := *n
		nn.b, nn.children = b, 0
		return &nn
	} else if n.b != b {
		// This node is from another block.
		// Ensure that it is given an identifier within that block
		// so it can be referenced from this one.
		if n.varname == "" {
			n.varname = n.b.genVarName()
		}
		return &node{b: b, expr: &ast.Identifier{Name: n.varname}}
	}
	return n
}

func (n *node) get() ast.Expression {
	if n.expr != nil {
		return n.expr
	}

	// We have never instantiated this node.
	// Perform that instantiation.
	n.expr = n.create()

	// If we have multiple children, we have
	// to wrap ourselves in a variable.
	// If a variable name has been specified,
	// try to use that.
	if n.children > 1 || n.varname != "" {
		if n.varname == "" {
			if _, ok := n.expr.(*ast.Identifier); ok {
				return n.expr
			}
			n.varname = n.b.genVarName()
		}
		n.b.vars[n.varname] = n.expr
		n.b.varnames = append(n.b.varnames, n.varname)
		n.expr = &ast.Identifier{Name: n.varname}
	}
	return n.expr
}
