package generate

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

func isArray(t *ast.Type) bool {
	return t.Elem != nil
}

func isPointer(t *ast.Type) bool {
	return !t.NonNull
}

func generateTypePrefix(t *ast.Type) string {
	typePrefix := ""
	if isPointer(t) {
		typePrefix += "*"
	}
	if isArray(t) {
		typePrefix += "[]"
		if isPointer(t.Elem) {
			typePrefix += "*"
		}
	}
	return typePrefix
}

func getFieldDefinitionTypeName(f *ast.FieldDefinition) string {
	return strings.Title(f.Type.Name()) + "Type"
}

func getFieldDefinitionName(f *ast.FieldDefinition) string {
	return strings.Title(f.Name)
}

func getBuildinTypeName(t *ast.Type) (string, bool) {
	buildinLookupKey := t.NamedType
	if t.Elem != nil {
		buildinLookupKey = t.Elem.NamedType
	}
	tn, ok := buildInTypeMap[buildinLookupKey]
	return tn, ok
}

func writeComment(buf *bytes.Buffer, description string) {
	if description == "" {
		return
	}
	fmt.Fprintf(buf, "// %s", description)
}

func stringLower(str string) string {
	return strings.ToLower(str[:1]) + str[1:]
}
