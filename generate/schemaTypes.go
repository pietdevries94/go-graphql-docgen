package generate

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pietdevries94/go-graphql-docgen/parser"
	"github.com/vektah/gqlparser/v2/ast"
	"log"
)

func GenerateSchemaTypes(buf *bytes.Buffer, parsed *parser.ParseResult, scalars map[string]string) {
	for _, td := range parsed.Schema.Types {
		if td.BuiltIn {
			continue
		}
		writeComment(buf, td.Description)
		fmt.Fprintf(buf, "type %sType ", strings.Title(td.Name))
		switch td.Kind {
		case ast.Object:
			generateObjectType(buf, td)
		case ast.Enum:
			generateEnum(buf, td)
		case ast.Scalar:
			generateScalar(buf, td, scalars)
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
		writeComment(buf, f.Description)
		typePrefix := generateTypePrefix(f.Type)
		typeName := getFieldDefinitionTypeName(f)
		if tn, ok := getBuildinTypeName(f.Type); ok {
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
		writeComment(buf, v.Description)
		fmt.Fprintf(buf, "%s%s %sType = \"%s\"\n", tn, strings.Title(v.Name), tn, v.Name)
	}
	fmt.Fprint(buf, ")")
}

func generateScalar(buf *bytes.Buffer, td *ast.Definition, scalars map[string]string) {
	tn, ok := scalars[td.Name]
	if !ok {
		log.Printf("WARNING: Unknown scalar %s, please define it in docgen.yml. Using interface{} as placeholder\n", td.Name)
		tn = "interface{}"
	}
	fmt.Fprint(buf, tn)
}
