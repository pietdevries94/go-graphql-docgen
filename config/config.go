package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SchemaFilename string `yaml:"schema"`
	QueriesFolder  string `yaml:"queries_folder"`
	Output         struct {
		File    string `yaml:"file"`
		Package string `yaml:"package"`
	} `yaml:"output"`
}

func LoadConfig() (*Config, error) {
	b, err := ioutil.ReadFile("./docgen.yml")
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}