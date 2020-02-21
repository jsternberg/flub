package flub

import "github.com/influxdata/flux/ast"

type Function struct {
	*Block
	n      *node
	params []*KeyValuePair
}

func (f *Function) node(b *Block) *node { return f.n.node(b) }
func (f *Function) get() ast.Expression { return f.n.get() }

func (f *File) Func(fn func(f *Function)) *Function {
	return f.body.Func(fn)
}

func (b *Block) Func(fn func(f *Function)) *Function {
	return &Function{
		n: b.newNode(func() ast.Expression {
			return Func(fn).get()
		}),
	}
}

func Func(fn func(f *Function)) *Function {
	f := &Function{
		Block: newBlock(),
	}
	f.n = &node{
		create: func() ast.Expression {
			return &ast.FunctionExpression{
				Params: asProperties(f.params),
				Body: func() ast.Node {
					stmts := f.Block.build()
					if len(stmts) == 1 {
						if e, ok := stmts[0].(*ast.ExpressionStatement); ok {
							return e.Expression
						}
					}
					return &ast.Block{
						Body: stmts,
					}
				}(),
			}
		},
	}
	fn(f)
	return f
}

func (f *Function) Arg(name string) Pipeline {
	n := &node{
		b:    f.Block,
		expr: &ast.Identifier{Name: name},
	}
	if !f.hasArg(name) {
		kvpair := &KeyValuePair{key: name}
		f.params = append(f.params, kvpair)
	}
	return Pipeline{n: n}
}

func (f *Function) ArgWithDefault(name string, def Node) Pipeline {
	n := &node{
		b:    f.Block,
		expr: &ast.Identifier{Name: name},
	}
	if !f.hasArg(name) {
		kvpair := &KeyValuePair{
			key:   name,
			value: def.node(f.n.b),
		}
		f.params = append(f.params, kvpair)
	}
	return Pipeline{n: n}
}

func (f *Function) hasArg(name string) bool {
	for _, kv := range f.params {
		if kv.key == name {
			return true
		}
	}
	return false
}
