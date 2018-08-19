package watch

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/uphy/elastic-watcher/pkg/actions"
	"github.com/uphy/elastic-watcher/pkg/condition"
	"github.com/uphy/elastic-watcher/pkg/context"
	"github.com/uphy/elastic-watcher/pkg/input"
	"github.com/uphy/elastic-watcher/pkg/transform"
	"github.com/uphy/elastic-watcher/pkg/trigger"
)

type (
	WatchConfig struct {
		Metadata  context.JSONObject    `json:"metadata,omitempty"`
		Trigger   *trigger.Trigger      `json:"trigger,omitempty"`
		Input     *input.Inputs         `json:"input,omitempty"`
		Condition *condition.Conditions `json:"condition,omitempty"`
		Transform *transform.Transforms `json:"transform,omitempty"`
		Actions   *actions.Actions      `json:"actions,omitempty"`
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
