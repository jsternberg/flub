package flub_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/influxdata/flux/ast"
	"github.com/influxdata/flux/parser"
	"github.com/jsternberg/flub"
)

func TestScript(t *testing.T) {
	for _, tt := range []struct {
		name string
		fn   func(f *flub.File)
		want string
	}{
		{
			name: "simple",
			fn: func(f *flub.File) {
				f.From("telegraf").
					Range(-time.Minute).
					Filter(flub.Func(func(fn *flub.Function) {
						r := fn.Arg("r")
						left := r.Get("_measurement").Eq(flub.String("cpu"))
						right := r.Get("_field").Eq(flub.String("usage_user"))
						left.And(right)
					}))
			},
			want: `package main

from(bucket: "telegraf")
	|> range(start: -1m)
	|> filter(fn: (r) => r._measurement == "cpu" and r._field == "usage_user")
`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			want := parser.ParseSource(tt.want).Files[0]
			f := flub.NewFile("main")
			tt.fn(f)
			got := f.Build()

			opts := []cmp.Option{
				cmpopts.IgnoreFields(ast.BaseNode{}, "Loc"),
				cmpopts.IgnoreFields(ast.File{}, "Metadata"),
			}
			if !cmp.Equal(want, got, opts...) {
				t.Log(ast.Format(got))
				t.Fatalf("unexpected output -want/+got:\n%s", cmp.Diff(want, got, opts...))
			}
		})
	}
}
