package pluginmanager

import "github.com/fredwangwang/bosh-plugin/monit"

func (p Manager) DisablePlugin(pluginName string) error {
	if err := monit.Stop(pluginName); err != nil {
		return err
	}

	states, err := p.ListPlugins()
	if err != nil {
		return err
	}

	for i, state := range states {
		if state.Name == pluginName {
			states[i].Enabled = false
		}
	}

	return WriteYamlStructToFile(states, p.configFilePath())
}
