package pluginmanager

import (
	"fmt"
	"github.com/fredwangwang/bosh-plugin/bpm"
	"github.com/fredwangwang/bosh-plugin/monit"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
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

	location := strings.Join(strings.Fields(info.Name), "")
	pluginPath := path.Join(p.Storage, location)

	cmd := exec.Command("unzip", filename, "-d", pluginPath)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to unzip plugin")
	}

	// create the folder/file structure required by BPM: /var/vcap/jobs/JOBNAME/config/bpm.yml
	jobPath := path.Join(p.Job, info.Name)
	configPath := path.Join(jobPath, "config")
	bpmConfigPath := path.Join(configPath, "bpm.yml")
	if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
		return err
	}

	var bpmConfig = bpm.Bpm{
		Processes: []bpm.Process{
			{
				Name:       info.Name,
				Executable: path.Join(jobPath, info.Command),
				Args:       info.Args,
				Env:        info.Env,
				Limits: map[string]string{
					"memory": info.Memory,
				},
			},
		},
	}

	if err := WriteYamlStructToFile(bpmConfig, bpmConfigPath); err != nil {
		return err
	}

	// create monit stub so monit are aware of the process
	if err := ioutil.WriteFile(
		path.Join(p.Monit, fmt.Sprintf("%s.monitrc", info.Name)),
		[]byte(monit.Monitrc(info.Name)),
		0644,
	); err != nil {
		return err
	}

	entrySrcPath := path.Join(pluginPath, info.Command)
	entryDstPath := path.Join(jobPath, info.Command)

	if err := copy.Copy(entrySrcPath, entryDstPath); err != nil {
		return err
	}

	if err := monit.Reload(); err != nil {
		return err
	}

	if err := monit.Start(info.Name); err != nil {
		return err
	}

	states = append(states, State{
		Name:        info.Name,
		Description: info.Description,
		Location:    location,
		Enabled:     true,
	})

	return WriteYamlStructToFile(states, p.configFilePath())
}
