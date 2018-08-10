package config

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
)

type (
	Config struct {
		Debug         bool          `json:"debug,omitempty"`
		Elasticsearch Elasticsearch `json:"elasticsearch"`
		Email         Email         `json:"email"`
	}
)

func LoadFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(f)
}

func Load(reader io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) Save(writer io.Writer) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func (c *Config) Validate() error {
	if strings.HasSuffix(c.Elasticsearch.URL, "/") {
		c.Elasticsearch.URL = c.Elasticsearch.URL[0 : len(c.Elasticsearch.URL)-1]
	}
	return nil
}
