package pluginmanager

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func (p Manager) ListPlugins() (States, error) {
	var states States
	var err error

	if configBytes, err := ioutil.ReadFile(p.configFilePath()); err != nil {
		return states, errors.Wrap(err, "failed to read config file")
	} else {
		err = yaml.Unmarshal(configBytes, &states)
	}

	return states, errors.Wrap(err, "failed to unmarshal config file")
}
