package flub

import (
	"regexp"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/ast"
)

func String(v string) Node {
	return &node{
		expr: &ast.StringLiteral{
			Value: v,
		},
	}
}

func Duration(d flux.Duration) Node {
	if d.IsNegative() {
		return &node{
			expr: &ast.UnaryExpression{
				Operator: ast.SubtractionOperator,
				Argument: &ast.DurationLiteral{
					Values: d.Mul(-1).AsValues(),
				},
			},
		}
	}
	return &node{
		expr: &ast.DurationLiteral{
			Values: d.AsValues(),
		},
	}
}

func Time(t flux.Time) Node {
	if t.IsRelative {
		return Duration(flux.ConvertDuration(t.Relative))
	}
	return &node{
		expr: &ast.DateTimeLiteral{
			Value: t.Absolute,
		},
	}
}

func Regex(regex *regexp.Regexp) Node {
	return &node{
		expr: &ast.RegexpLiteral{
			Value: regex,
		},
	}
}

func True() Node {
	return &node{
		expr: &ast.BooleanLiteral{
			Value: true,
		},
	}
}

func False() Node {
	return &node{
		expr: &ast.BooleanLiteral{
			Value: false,
		},
	}
}

func Strings(arr []string) Node {
	elements := make([]ast.Expression, len(arr))
	for i := range arr {
		elements[i] = &ast.StringLiteral{Value: arr[i]}
	}
	return &node{
		expr: &ast.ArrayExpression{
			Elements: elements,
		},
	}
}
