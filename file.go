package flub

import "github.com/influxdata/flux/ast"

type File struct {
	pkgname string
	body    *Block
}

func NewFile(pkgname string) *File {
	return &File{
		pkgname: pkgname,
		body:    newBlock(),
	}
}

func (f *File) Call(name string, args ...*KeyValuePair) Pipeline {
	return f.body.Call(name, args...)
}

func (f *File) From(bucket string) Pipeline {
	return f.body.From(bucket)
}

// Build will construct the AST file.
func (f *File) Build() *ast.File {
	return &ast.File{
		Package: &ast.PackageClause{
			Name: &ast.Identifier{Name: f.pkgname},
		},
		Body: f.body.build(),
	}
}

func (f *File) Format() string {
	return ast.Format(f.Build())
}
