package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/pietdevries94/go-graphql-docgen/config"
	"github.com/pietdevries94/go-graphql-docgen/generate"
	"github.com/pietdevries94/go-graphql-docgen/parser"
)

// const endpoint = "https://favware.tech/api"

func Execute() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	ensureFolder(cfg.Output.Folder)

	parsed := parser.MustParse(cfg)

	schemaTypesBuf := newFileBuffer(cfg.Output.Package)
	generate.GenerateSchemaTypes(schemaTypesBuf, parsed, cfg.Scalars)
	writeFile(cfg.Output.Folder, "schemaTypes.go", schemaTypesBuf)

	queriesBuf := newFileBuffer(cfg.Output.Package)
	generate.GenerateQueries(queriesBuf, parsed)
	writeFile(cfg.Output.Folder, "queries.go", queriesBuf)

	cmd := exec.Command(`gofmt`, `-w`, cfg.Output.Folder)
	err = cmd.Run()
	if err != nil {
		log.Fatal("FATAL: gofmt failed")
	}
}

func newFileBuffer(packageName string) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	fmt.Fprint(buf, "// Code generated by github.com/pietdevries94/go-graphql-docgen, DO NOT EDIT.\n\n")
	fmt.Fprintf(buf, "package %s\n\n", packageName)
	return buf
}

func writeFile(folder, name string, buf *bytes.Buffer) {
	fn := path.Join(folder, name)
	err := ioutil.WriteFile(fn, buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func ensureFolder(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
}
