package generate

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

func createBaseClient(buf *bytes.Buffer) {
	fmt.Fprint(buf, `
	import (
		"github.com/machinebox/graphql"

		"context"
	)

	type Client struct {
		c *graphql.Client
	}

	type RequestOptions struct {}
	
	func NewClient(endpoint string, opts ...graphql.ClientOption) *Client {
		c := graphql.NewClient(endpoint, opts...)
		return &Client{c: c}
	}

	func (client *Client) run(req *graphql.Request, options *RequestOptions, respData interface{}) error {
		ctx := context.Background()
		return client.c.Run(ctx, req, respData)
	}
`)
}

func generateClientFunction(buf *bytes.Buffer, op *ast.OperationDefinition, name string) {
	docName := strings.ToLower(name[0:1]) + name[1:] + "Document"
	fmt.Fprintf(buf, "const %s = `%s`\n", docName, op.Position.Src.Input)

	fmt.Fprintf(buf, "func (c *Client) %s(", name)

	setVarsBuf := bytes.NewBuffer(nil)
	if len(op.VariableDefinitions) > 0 {
		for _, varDef := range op.VariableDefinitions {
			typePrefix := ""
			if isPointer(varDef.Type) {
				typePrefix = "*"
			}
			tn := strings.Title(varDef.Type.Name()) + "Type"
			if bt, ok := getBuildinTypeName(varDef.Type); ok {
				tn = bt
			}
			fmt.Fprintf(buf, "%s %s%s, ", varDef.Variable, typePrefix, tn)

			// set the variable in request
			fmt.Fprintf(setVarsBuf, "req.Var(\"%s\", %s)\n", varDef.Variable, varDef.Variable)
		}
	}
	fmt.Fprintf(buf, "options *RequestOptions) (*%sQueryResult, error) {\n", name)

	fmt.Fprintf(buf, "req := graphql.NewRequest(%s)\n", docName)
	_, err := setVarsBuf.WriteTo(buf)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(buf, "respData := &%sQueryResult{}\n", name)
	fmt.Fprint(buf, `
	err := c.run(req, options, respData)
	return respData, err
	`)
	fmt.Fprint(buf, "}\n\n")
}
