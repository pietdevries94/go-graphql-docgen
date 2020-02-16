package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/pietdevries94/go-graphql-docgen/config"
	"github.com/pietdevries94/go-graphql-docgen/generate"
	"github.com/pietdevries94/go-graphql-docgen/parser"
)

// const endpoint = "https://countries.trevorblades.com/"

func Execute() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	parsed := parser.MustParse(cfg)
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "package %s\n\n", cfg.Output.Package)
	res := generate.GenerateQueries(buf, parsed)

	err = ioutil.WriteFile(cfg.Output.File, res.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(`gofmt`, `-w`, cfg.Output.File)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
