package parser

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/pietdevries94/go-graphql-docgen/config"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/gqlerror"
	"github.com/vektah/gqlparser/parser"
	"github.com/vektah/gqlparser/validator"
)

type ParseResult struct {
	Schema  *ast.Schema
	Queries []*ast.QueryDocument
}

func MustParse(cfg *config.Config) *ParseResult {
	res, err := Parse(cfg)
	if err != nil {
		panic(err)
	}
	return res
}

func Parse(cfg *config.Config) (*ParseResult, error) {
	res := &ParseResult{}
	var err error
	res.Schema, err = getSchema(cfg.SchemaFilename)
	if err != nil {
		return nil, err
	}

	qSources, err := getQuerySources(cfg.QueriesFolder)
	if err != nil {
		return nil, err
	}

	res.Queries = make([]*ast.QueryDocument, len(qSources))
	for i, qSrc := range qSources {
		doc, err := loadQuery(res.Schema, qSrc)
		if err != nil {
			return nil, err
		}
		res.Queries[i] = doc
	}
	return res, nil
}

func loadQuery(schema *ast.Schema, src *ast.Source) (*ast.QueryDocument, gqlerror.List) {
	query, err := parser.ParseQuery(src)
	if err != nil {
		return nil, gqlerror.List{err}
	}
	errs := validator.Validate(schema, query)
	if errs != nil {
		return nil, errs
	}

	return query, nil
}

func getSchema(filePath string) (*ast.Schema, error) {
	src, err := getSchemaSource()
	if err != nil {
		return nil, err
	}
	schema, gqlErr := gqlparser.LoadSchema(src)
	if gqlErr != nil {
		return nil, gqlErr
	}
	return schema, nil
}

func getSchemaSource() (*ast.Source, error) {
	f, err := os.Open("./testdata/schema.graphql")
	if err != nil {
		return nil, err
	}

	input, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	s := &ast.Source{
		Name:  f.Name(),
		Input: string(input),
	}
	return s, nil
}

func getQuerySources(folderPath string) ([]*ast.Source, error) {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	res := []*ast.Source{}
	for _, file := range files {
		fPath := path.Join(folderPath, file.Name())
		if file.IsDir() {
			content, err := getQuerySources(fPath)
			if err != nil {
				return nil, err
			}
			res = append(res, content...)
			continue
		}

		ext := path.Ext(file.Name())
		if ext != ".graphql" && ext != ".gql" {
			continue
		}

		input, err := ioutil.ReadFile(fPath)
		if err != nil {
			return nil, err
		}
		res = append(res, &ast.Source{
			Name:  file.Name(),
			Input: string(input),
		})
	}
	return res, nil
}
