package generate

import "github.com/vektah/gqlparser/v2/ast"

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
