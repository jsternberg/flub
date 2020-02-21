package flub

import "github.com/influxdata/flux/ast"

type KeyValuePair struct {
	key   string
	value *node
}

func KV(key string, value Node) *KeyValuePair {
	n := value.node(nil)
	n.children++
	return &KeyValuePair{
		key:   key,
		value: n,
	}
}

func asProperties(pairs []*KeyValuePair) []*ast.Property {
	properties := make([]*ast.Property, 0, len(pairs))
	for _, kv := range pairs {
		var value ast.Expression
		if kv.value != nil {
			value = kv.value.get()
		}
		properties = append(properties, &ast.Property{
			Key:   &ast.Identifier{Name: kv.key},
			Value: value,
		})
	}
	return properties
}
