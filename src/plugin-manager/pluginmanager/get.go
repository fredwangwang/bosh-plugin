package pluginmanager

import "fmt"

func (p Manager) GetPlugin(pluginName string) (State, error) {
	states, err := p.ListPlugins()
	if err != nil {
		return State{}, err
	}

	for _, state := range states {
		if state.Name == pluginName {
			return state, nil
		}
	}

	return State{}, fmt.Errorf("failed to find %s", pluginName)
}
