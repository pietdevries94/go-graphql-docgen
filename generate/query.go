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

func GenerateQueries(buf *bytes.Buffer, parsed *parser.ParseResult, generateClient bool) {
	if generateClient {
		createBaseClient(buf)
	}

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

			generateStruct(buf, frag.SelectionSet, fmt.Sprintf("%sFragment", strings.Title(name)))
		}
		for _, op := range query.Operations {
			name := op.Name
			if name == "" {
				name = nameFromFileName(query.Position.Src.Name)
			}
			name = strings.Title(name)

			generateStruct(buf, op.SelectionSet, fmt.Sprintf("%sQueryResult", name))

			if generateClient {
				generateClientFunction(buf, op, name)
			}
		}
	}
}

type structGenerator struct {
	bufs []*bytes.Buffer
}

func generateStruct(buf *bytes.Buffer, set ast.SelectionSet, name string) {
	sg := structGenerator{[]*bytes.Buffer{}}
	sg.generateStruct(set, name)
	for _, b := range sg.bufs {
		_, err := buf.ReadFrom(b)
		if err != nil {
			panic(err)
		}
	}
}

func (sg *structGenerator) generateStruct(set ast.SelectionSet, name string) {
	buf := bytes.NewBufferString("type ")
	sg.bufs = append(sg.bufs, buf)
	fmt.Fprintf(buf, "%s struct {\n", name)
	for _, sel := range set {
		sg.parseSelection(buf, sel, name)
	}
	buf.WriteString("}\n\n")
}

func (sg *structGenerator) parseSelection(buf *bytes.Buffer, sel ast.Selection, name string) {
	if f, ok := sel.(*ast.Field); ok {
		sg.writeField(buf, f, name)
	} else if frag, ok := sel.(*ast.FragmentSpread); ok {
		sg.writeFragment(buf, frag)
	}
}

func (sg *structGenerator) writeField(buf *bytes.Buffer, f *ast.Field, name string) {
	typePrefix := generateTypePrefix(f.Definition.Type)

	if isSingleFragment(f.SelectionSet) {
		fmt.Fprintf(buf, "%s %s", strings.Title(f.Alias), typePrefix)
		sg.writeFragment(buf, f.SelectionSet[0].(*ast.FragmentSpread))
		return
	}

	if f.SelectionSet != nil {
		tn := name + strings.Title(f.Alias)
		fmt.Fprintf(buf, "%s %s%s\n\n", strings.Title(f.Alias), typePrefix, tn)

		sg.generateStruct(f.SelectionSet, tn)
		return
	}

	if bt, ok := getBuildinTypeName(f.Definition.Type); ok {
		fmt.Fprintf(buf, "%s %s%s\n", strings.Title(f.Alias), typePrefix, bt)
	}
}

func (sg *structGenerator) writeFragment(buf *bytes.Buffer, f *ast.FragmentSpread) {
	fmt.Fprintf(buf, "*%sFragment\n", strings.Title(f.Name))
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
