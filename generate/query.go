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
			buf.WriteString("}\n")
		}
	}

	return buf
}

func writeField(buf *bytes.Buffer, f *ast.Field) {
	if f.SelectionSet != nil {
		fmt.Fprintf(buf, "%s struct {\n", strings.Title(f.Alias))
		for _, sel := range f.SelectionSet {
			if f, ok := sel.(*ast.Field); ok {
				writeField(buf, f)
			}
		}
		fmt.Fprintf(buf, "} `%s`\n", getTags(f))
		return
	}

	if bt, ok := baseTypeMap[f.Definition.Type.Name()]; ok {
		if !f.Definition.Type.NonNull {
			bt = "*" + bt
		}
		fmt.Fprintf(buf, "%s %s `%s`\n", strings.Title(f.Alias), bt, getTags(f))
		return
	}
}

func nameFromFileName(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func getTags(f *ast.Field) string {
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
			t += " "
		}
		t += k + ":'" + v + "'"
	}
	return t + `"`
}
