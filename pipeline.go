package flub

import (
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/ast"
)

type Pipeline struct {
	n *node
}

func (p Pipeline) As(name string) Pipeline {
	if p.n.varname != "" {
		panic("variable name redeclared")
	}
	p.n.varname = name
	return Pipeline{n: p.n}
}

func (p Pipeline) Call(name string, args ...*KeyValuePair) Pipeline {
	arguments := p.readArgs(args)
	n := p.n.b.newNode(func() ast.Expression {
		return &ast.PipeExpression{
			Argument: p.n.get(),
			Call:     Call(name, arguments...),
		}
	})
	p.n.children++
	return Pipeline{n: n}
}

func (p Pipeline) Invoke(args ...*KeyValuePair) Pipeline {
	p.n.children++
	arguments := p.readArgs(args)
	return Pipeline{n: p.n.b.newNode(func() ast.Expression {
		return Invoke(p.n.get(), arguments...)
	})}
}

func (p Pipeline) readArgs(args []*KeyValuePair) []*KeyValuePair {
	newArgs := make([]*KeyValuePair, len(args))
	for i, arg := range args {
		newArgs[i] = &KeyValuePair{
			key:   arg.key,
			value: arg.value.node(p.n.b),
		}
	}
	return newArgs
}

func (p Pipeline) Get(key string) Pipeline {
	p.n.children++
	return Pipeline{n: p.n.b.newNode(func() ast.Expression {
		return &ast.MemberExpression{
			Object:   p.n.get(),
			Property: &ast.Identifier{Name: key},
		}
	})}
}

func (p Pipeline) And(other Node) Pipeline { return p.logicalExpr(ast.AndOperator, other) }
func (p Pipeline) Or(other Node) Pipeline  { return p.logicalExpr(ast.OrOperator, other) }

func (p Pipeline) logicalExpr(op ast.LogicalOperatorKind, other Node) Pipeline {
	p.n.children++
	n := other.node(p.n.b)
	n.children++
	return Pipeline{n: p.n.b.newNode(func() ast.Expression {
		return &ast.LogicalExpression{
			Operator: op,
			Left:     p.n.get(),
			Right:    n.get(),
		}
	})}
}

func (p Pipeline) Eq(other Node) Pipeline {
	return p.binaryExpr(ast.EqualOperator, other)
}
func (p Pipeline) RegexEq(other Node) Pipeline {
	return p.binaryExpr(ast.RegexpMatchOperator, other)
}

func (p Pipeline) binaryExpr(op ast.OperatorKind, other Node) Pipeline {
	p.n.children++
	n := other.node(p.n.b)
	n.children++
	return Pipeline{n: p.n.b.newNode(func() ast.Expression {
		return &ast.BinaryExpression{
			Operator: op,
			Left:     p.n.get(),
			Right:    n.get(),
		}
	})}
}

func (p Pipeline) Range(start time.Duration) Pipeline {
	return p.Call("range",
		KV("start", Duration(flux.ConvertDuration(start))),
	)
}

func (p Pipeline) Filter(fn *Function) Pipeline {
	return p.Call("filter",
		KV("fn", fn),
	)
}

func (p Pipeline) Group(columns []string) Pipeline {
	return p.Call("group",
		KV("columns", Strings(columns)),
	)
}

func (p Pipeline) Sort(columns []string) Pipeline {
	return p.Call("sort",
		KV("columns", Strings(columns)),
	)
}

func (p Pipeline) Keep(columns []string) Pipeline {
	return p.Call("keep",
		KV("columns", Strings(columns)),
	)
}

func (p Pipeline) Yield(name string) Pipeline {
	return p.Call("yield",
		KV("name", String(name)),
	)
}

func (p Pipeline) node(b *Block) *node { return p.n.node(b) }
func (p Pipeline) get() ast.Expression { return p.n.get() }

func Call(name string, args ...*KeyValuePair) *ast.CallExpression {
	callee := &ast.Identifier{Name: name}
	return Invoke(callee, args...)
}

func Invoke(callee ast.Expression, args ...*KeyValuePair) *ast.CallExpression {
	var arguments []ast.Expression
	if len(args) > 0 {
		arguments = []ast.Expression{
			&ast.ObjectExpression{Properties: asProperties(args)},
		}
	}
	return &ast.CallExpression{
		Callee:    callee,
		Arguments: arguments,
	}
}
