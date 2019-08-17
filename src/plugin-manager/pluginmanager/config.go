package pluginmanager

import (
	"net/url"
)

func (p Manager) ConfigPlugin(pluginName string, query url.Values) error {
	stats, err := p.ListPlugins()
	if err != nil {
		return err
	}
	pstat, err := GetPState(stats, pluginName)
	if err != nil {
		return err
	}

	for queryKey := range query {
		pstat.PendingEnv[queryKey] = query.Get(queryKey)
	}

	return WriteYamlStructToFile(stats, p.configFilePath())
}
