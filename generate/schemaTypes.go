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
		case ast.Enum:
			generateEnum(buf, td)
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
		typeName := getFieldDefinitionTypeName(f) + "Type"
		if tn, ok := buildInTypeMap[f.Type.NamedType]; ok {
			typeName = tn
		}
		fmt.Fprintf(buf, "%s %s%s\n", getFieldDefinitionName(f), typePrefix, typeName)
	}
	fmt.Fprint(buf, "}")
}

func generateEnum(buf *bytes.Buffer, td *ast.Definition) {
	fmt.Fprintf(buf, "string\n")

	tn := strings.Title(td.Name)
	fmt.Fprint(buf, "\nconst (")
	for _, v := range td.EnumValues {
		fmt.Fprintf(buf, "%s%s %sType = \"%s\"\n", strings.Title(v.Name), tn, tn, v.Name)
	}
	fmt.Fprint(buf, ")")
}
