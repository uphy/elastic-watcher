package watcher

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/uphy/elastic-watcher/watcher/actions"
	"github.com/uphy/elastic-watcher/watcher/condition"
	"github.com/uphy/elastic-watcher/watcher/input"
	"github.com/uphy/elastic-watcher/watcher/transform"
)

type (
	WatchConfig struct {
		Metadata  map[string]interface{} `json:"metadata"`
		Trigger   Trigger                `json:"trigger"`
		Input     input.Input            `json:"input"`
		Condition condition.Conditions   `json:"condition"`
		Transform *transform.Transform   `json:"transform"`
		Actions   actions.Actions        `json:"actions"`
	}
	Trigger struct {
		Schedule *Schedule `json:"schedule"`
	}
)

func LoadFile(file string) (*WatchConfig, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(f)
}

func Load(reader io.Reader) (*WatchConfig, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var c WatchConfig
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *WatchConfig) Save(writer io.Writer) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}
