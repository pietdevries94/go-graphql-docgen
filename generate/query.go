package generate

import (
	"bytes"
	"fmt"
	"github.com/vektah/gqlparser/ast"
	"path"
	"strings"

	"github.com/pietdevries94/go-graphql-docgen/parser"
)

var generatedFragments = []string{}

func GenerateQueries(buf *bytes.Buffer, parsed *parser.ParseResult) {
	for _, query := range parsed.Queries {
	fragmentLoop:
		for _, frag := range query.Fragments {
			name := frag.Name
			if name == "" {
				name = nameFromFileName(query.Position.Src.Name)
			}
			for _, fn := range generatedFragments {
				if fn == name {
					continue fragmentLoop
				}
			}
			generatedFragments = append(generatedFragments, name)
			fmt.Fprintf(buf, "type %sFragment ", strings.Title(name))
			generateStruct(buf, frag.SelectionSet)
		}
		for _, op := range query.Operations {
			name := op.Name
			if name == "" {
				name = nameFromFileName(query.Position.Src.Name)
			}
			fmt.Fprintf(buf, "type %sQueryResult ", strings.Title(name))
			generateStruct(buf, op.SelectionSet)
		}
	}
}

func generateStruct(buf *bytes.Buffer, set ast.SelectionSet) {
	buf.WriteString("struct {\n")
	for _, sel := range set {
		if f, ok := sel.(*ast.Field); ok {
			writeField(buf, f)
		}
	}
	buf.WriteString("}\n\n")
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
			} else if frag, ok := sel.(*ast.FragmentSpread); ok {
				writeFragment(buf, frag)
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

func writeFragment(buf *bytes.Buffer, f *ast.FragmentSpread) {
	fmt.Fprintf(buf, "*%sFragment", f.Name)
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
