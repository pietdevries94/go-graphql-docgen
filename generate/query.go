package generate

import (
	"bytes"
	"fmt"
	"github.com/vektah/gqlparser/ast"
	"path"
	"strings"

	"github.com/pietdevries94/go-graphql-docgen/parser"
)

func GenerateQueries(buf *bytes.Buffer, parsed *parser.ParseResult) *bytes.Buffer {
	for _, query := range parsed.Queries {
		for _, op := range query.Operations {
			name := op.Name
			if name == "" {
				name = nameFromFileName(query.Position.Src.Name)
			}
			fmt.Fprintf(buf, "type %sQueryResult struct {\n", strings.Title(name))
			for _, sel := range op.SelectionSet {
				if f, ok := sel.(*ast.Field); ok {
					writeField(buf, f)
				}
			}
			buf.WriteString("}\n\n")
		}
	}

	return buf
}

func writeField(buf *bytes.Buffer, f *ast.Field) {
	typePrefix := ""
	if isPointer(f) {
		typePrefix = "*"
	}
	if isArray(f) {
		typePrefix += "[]"
	}

	if f.SelectionSet != nil {
		fmt.Fprintf(buf, "%s %sstruct {\n", strings.Title(f.Alias), typePrefix)
		for _, sel := range f.SelectionSet {
			if f, ok := sel.(*ast.Field); ok {
				writeField(buf, f)
			}
		}
		fmt.Fprintf(buf, "} `%s`\n", getFieldTags(f))
		return
	}

	if bt, ok := baseTypeMap[f.Definition.Type.Name()]; ok {
		fmt.Fprintf(buf, "%s %s%s `%s`\n", strings.Title(f.Alias), typePrefix, bt, getFieldTags(f))
		return
	}
}

func nameFromFileName(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func getFieldTags(f *ast.Field) string {
	tm := map[string]string{
		`type`: f.Definition.Type.Name(),
	}
	if f.Name != f.Alias {
		tm[`name`] = f.Name
	}
	t := `docgen:"`
	first := true
	for k, v := range tm {
		if first {
			first = false
		} else {
			t += ", "
		}
		t += k + ":'" + v + "'"
	}
	return t + `"`
}

func isArray(f *ast.Field) bool {
	return f.Definition.Type.Elem != nil
}

func isPointer(f *ast.Field) bool {
	return !f.Definition.Type.NonNull
}
