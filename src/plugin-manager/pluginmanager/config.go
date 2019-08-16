package pluginmanager

import (
	"net/url"
)

func (p Manager) ConfigPlugin(pluginName string, query url.Values) error {
	var stat State
	var idx int

	stats, err := p.ListPlugins()
	if err != nil {
		return err
	}

	for idx, stat = range stats {
		if stat.Name == pluginName {
			break
		}
	}

	for queryKey := range query {
		stat.PendingEnv[queryKey] = query.Get(queryKey)
	}
	stats[idx] = stat

	return WriteYamlStructToFile(stats, p.configFilePath())
}
