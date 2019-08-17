package pluginmanager

import (
	"github.com/fredwangwang/bosh-plugin/bpm"
	"github.com/fredwangwang/bosh-plugin/monit"
)

func (p Manager) EnablePlugin(pluginName string) error {
	states, err := p.ListPlugins()
	if err != nil {
		return err
	}
	pstat, err := GetPState(states, pluginName)
	if err != nil {
		return err
	}

	// already enabled, nothing to do
	if pstat.Enabled == true {
		return nil
	}

	pstat.Enabled = true
	pstat.AdditionalEnv = mergeMaps(pstat.AdditionalEnv, pstat.PendingEnv)
	pstat.PendingEnv = map[string]string{}

	var bpmConfig bpm.Bpm
	if err := ReadYamlStructFromFile(p.pluginBPMPath(pluginName), &bpmConfig); err != nil {
		return err
	}
	bpmConfig.Processes[0].Env = mergeMaps(pstat.Env, pstat.AdditionalEnv)

	if err := WriteYamlStructToFile(bpmConfig, p.pluginBPMPath(pluginName)); err != nil {
		return err
	}

	if err := monit.Start(pluginName); err != nil {
		return err
	}

	return WriteYamlStructToFile(states, p.configFilePath())
}

func mergeMaps(maps ...map[string]string) map[string]string {
	resMap := map[string]string{}
	for _, m := range maps {
		for k, v := range m {
			if v == "" {
				continue
			}
			resMap[k] = v
		}
	}
	return resMap
}
