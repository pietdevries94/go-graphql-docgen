package generate

import (
	"bytes"
	"fmt"
	"github.com/vektah/gqlparser/v2/ast"
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

			if len(op.VariableDefinitions) > 0 {
				generateVariablesStruct(buf, op, name)
			}
		}
	}
}

func generateVariablesStruct(buf *bytes.Buffer, op *ast.OperationDefinition, name string) {
	fmt.Fprintf(buf, "type %sQueryVariables struct {\n", strings.Title(name))
	for _, varDef := range op.VariableDefinitions {
		typePrefix := ""
		if isPointer(varDef.Type) {
			typePrefix = "*"
		}
		if bt, ok := buildInTypeMap[varDef.Type.Name()]; ok {
			fmt.Fprintf(buf, "%s %s%s\n", name, typePrefix, bt)
		}
	}
	buf.WriteString("}\n\n")
}

func generateStruct(buf *bytes.Buffer, set ast.SelectionSet) {
	buf.WriteString("struct {\n")
	for _, sel := range set {
		parseSelection(buf, sel)
	}
	buf.WriteString("}\n\n")
}

func parseSelection(buf *bytes.Buffer, sel ast.Selection) {
	if f, ok := sel.(*ast.Field); ok {
		writeField(buf, f)
	} else if frag, ok := sel.(*ast.FragmentSpread); ok {
		writeFragment(buf, frag)
	}
}

func writeField(buf *bytes.Buffer, f *ast.Field) {
	typePrefix := generateTypePrefix(f.Definition.Type)

	if isSingleFragment(f.SelectionSet) {
		fmt.Fprintf(buf, "%s %s", strings.Title(f.Alias), typePrefix)
		writeFragment(buf, f.SelectionSet[0].(*ast.FragmentSpread))
		return
	}

	if f.SelectionSet != nil {
		fmt.Fprintf(buf, "%s %sstruct {\n", strings.Title(f.Alias), typePrefix)
		for _, sel := range f.SelectionSet {
			parseSelection(buf, sel)
		}
		fmt.Fprintf(buf, "} `%s`\n", getFieldTags(f))
		return
	}

	if bt, ok := buildInTypeMap[f.Definition.Type.Name()]; ok {
		fmt.Fprintf(buf, "%s %s%s `%s`\n", strings.Title(f.Alias), typePrefix, bt, getFieldTags(f))
		return
	}
}

func writeFragment(buf *bytes.Buffer, f *ast.FragmentSpread) {
	fmt.Fprintf(buf, "*%sFragment\n", f.Name)
}

func nameFromFileName(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func isSingleFragment(sel ast.SelectionSet) bool {
	if sel != nil && len(sel) == 1 {
		_, ok := (sel)[0].(*ast.FragmentSpread)
		return ok
	}
	return false
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
