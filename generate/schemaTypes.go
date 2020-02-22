package generate

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pietdevries94/go-graphql-docgen/parser"
	"github.com/vektah/gqlparser/v2/ast"
	"log"
)

func GenerateSchemaTypes(buf *bytes.Buffer, parsed *parser.ParseResult) {
	for _, td := range parsed.Schema.Types {
		if td.BuiltIn {
			continue
		}
		fmt.Fprintf(buf, "type %sType ", strings.Title(td.Name))
		switch td.Kind {
		case ast.Object:
			generateObjectType(buf, td)
		default:
			log.Printf("TODO: %s\n", td.Kind)
			fmt.Fprint(buf, "interface{}")
		}
		fmt.Fprint(buf, "\n\n")
	}
}

func generateObjectType(buf *bytes.Buffer, td *ast.Definition) {
	fmt.Fprint(buf, "struct{\n")
	for _, f := range td.Fields {
		typePrefix := generateTypePrefix(f.Type)
		typeName := f.Type.Name() + "Type"
		if tn, ok := buildInTypeMap[f.Type.NamedType]; ok {
			typeName = tn
		}
		fmt.Fprintf(buf, "%s %s%s\n", strings.Title(f.Name), typePrefix, typeName)
	}
	fmt.Fprint(buf, "}")
}
