package pluginmanager

import (
	"fmt"
	"github.com/fredwangwang/bosh-plugin/monit"
	"github.com/pkg/errors"
	"os"
)

func (p Manager) DeletePlugin(pluginName string) error {
	if _, err := os.Stat(p.pluginStorePath(pluginName)); os.IsNotExist(err) {
		return fmt.Errorf("plugin %s does not exist", pluginName)
	}

	if err := monit.Stop(pluginName); err != nil {
		return err
	}

	if err := os.RemoveAll(p.pluginStorePath(pluginName)); err != nil {
		return err
	}

	if err := os.RemoveAll(p.pluginJobPath(pluginName)); err != nil {
		return err
	}

	if err := os.RemoveAll(p.pluginMonitPath(pluginName)); err != nil {
		return err
	}

	if err := monit.Reload(); err != nil {
		return err
	}

	states, err := p.ListPlugins()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve existing plugins")
	}

	var pluginIdx int
	for pluginIdx = 0; pluginIdx < len(states); pluginIdx++ {
		if states[pluginIdx].Name == pluginName {
			break
		}
	}
	copy(states[pluginIdx:], states[pluginIdx+1:])
	states[len(states)-1] = State{} // or the zero value of T
	states = states[:len(states)-1]

	return WriteYamlStructToFile(states, p.configFilePath())
}
