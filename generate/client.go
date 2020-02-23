package generate

import (
	"bytes"
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
)

func createBaseClient(buf *bytes.Buffer) {
	fmt.Fprint(buf, `
	import "github.com/machinebox/graphql"

	type Client struct {
		c *graphql.Client
	}
	
	func NewClient(endpoint string, opts ...graphql.ClientOption) *Client {
		c := graphql.NewClient(endpoint, opts...)
		return &Client{c: c}
	}
`)
}

func generateClientFunction(buf *bytes.Buffer, op *ast.OperationDefinition, name string) {
	funcArgs := ""
	if len(op.VariableDefinitions) > 0 {
		funcArgs = fmt.Sprintf("vars %sQueryVariables", name)
	}
	fmt.Fprintf(buf, "func (c *Client) %s(%s) *%sQueryResult {\npanic(`TODO`)\n}\n\n", name, funcArgs, name)
}
