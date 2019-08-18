package pluginmanager

import (
	"github.com/fredwangwang/bosh-plugin/bpm"
	"github.com/fredwangwang/bosh-plugin/monit"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path"
)

func (p Manager) AddPlugin(filename string) error {
	var err error

	// TODO: handle dup with existing ones
	states, err := p.ListPlugins()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve existing plugins")
	}

	infos, err := getPluginInfo(filename)
	if err != nil {
		return err
	}

	info := infos.Applications[0]

	pluginPath := p.pluginStorePath(info.Name)

	cmd := exec.Command("unzip", filename, "-d", pluginPath)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to unzip plugin")
	}

	// create the folder/file structure required by BPM in the plugin storage path
	// this will be copy over to the correct job path in the following step
	configPathInStore := path.Join(pluginPath, "config")
	bpmConfigPathInStore := path.Join(configPathInStore, "bpm.yml")
	if err := os.MkdirAll(configPathInStore, os.ModePerm); err != nil {
		return err
	}

	var bpmConfig = bpm.Bpm{
		Processes: []bpm.Process{
			{
				Name:       info.Name,
				Executable: path.Join(p.pluginJobPath(info.Name), info.Command),
				Args:       info.Args,
				Env:        p.addMetronEnv(info.Name, info.Env),
				Limits: map[string]string{
					"memory": info.Memory,
				},
			},
		},
	}

	if err := WriteYamlStructToFile(bpmConfig, bpmConfigPathInStore); err != nil {
		return err
	}

	if err := p.copyMetronFilesFor(info.Name, configPathInStore); err != nil {
		return err
	}

	if err := p.copyPluginFromStorageToJob(info.Name); err != nil {
		return err
	}

	// create monit stub so monit are aware of the process
	if err := monit.CreateMonitrcFor(info.Name, p.pluginMonitPath(info.Name)); err != nil {
		return err
	}

	if err := monit.Reload(); err != nil {
		return err
	}

	if err := monit.Start(info.Name); err != nil {
		return err
	}

	states = append(states, State{
		Name:          info.Name,
		Description:   info.Description,
		Location:      info.Name,
		Enabled:       true,
		Env:           info.Env,
		Arg:           info.Args,
		AdditionalEnv: p.addMetronEnv(info.Name, nil),
		PendingEnv:    map[string]string{},
	})

	return WriteYamlStructToFile(states, p.configFilePath())
}

func (p Manager) copyPluginFromStorageToJob(pluginName string) error {

	if err := copy.Copy(p.pluginStorePath(pluginName), p.pluginJobPath(pluginName)); err != nil {
		return err
	}

	if err := exec.Command("chgrp", "-R", "vcap", p.pluginJobPath(pluginName)).Start(); err != nil {
		return err
	}
	return nil
}

func (p Manager) copyMetronFilesFor(pluginName string, dstConfigPath string) error {
	srcConfigPath := path.Join(p.Job, "plugin-manager", "config")

	for _, file := range METRON_FILES {
		if err := copy.Copy(path.Join(srcConfigPath, file), path.Join(dstConfigPath, file)); err != nil {
			return err
		}
	}
	return nil
}

func (p Manager) addMetronEnv(pluginName string, env map[string]string) map[string]string {
	res := map[string]string{}

	for k, v := range env {
		res[k] = v
	}

	dstConfigPath := path.Join(p.pluginJobPath(pluginName), "config")
	for k, v := range METRON_FILES {
		res[k] = path.Join(dstConfigPath, v)
	}

	return res
}

var METRON_FILES = map[string]string{
	"METRON_CA_CERT_PATH": "metron_ca_cert.pem",
	"METRON_CERT_PATH":    "metron_cert.pem",
	"METRON_KEY_PATH":     "metron_cert.key",
}
