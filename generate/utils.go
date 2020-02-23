package generate

import (
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
	return strings.Title(f.Type.Name())
}

func getFieldDefinitionName(f *ast.FieldDefinition) string {
	return strings.Title(f.Name)
}
