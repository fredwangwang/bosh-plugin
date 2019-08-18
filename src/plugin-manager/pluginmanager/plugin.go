package pluginmanager

import (
	"archive/zip"
	"fmt"
	"github.com/fredwangwang/bosh-plugin/monit"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
)

const PLUGIN_INFO_FILE = "plugin.yml"

type Manager struct {
	Job        string
	Monit      string
	Storage    string
	ConfigFile string
}

type State struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Location    string `yaml:"location"`
	Enabled     bool   `yaml:"enabled"`

	Env           map[string]string `yaml:"env"`
	Arg           []string          `yaml:"arg"`
	AdditionalEnv map[string]string `yaml:"additional-env"`
	PendingEnv    map[string]string `yaml:"pending-env"`
}

type States []State

func GetPluginManager(job string, monit string, storage string, configFile string) (Manager, error) {
	pm := Manager{
		Job:        job,
		Monit:      monit,
		Storage:    storage,
		ConfigFile: configFile,
	}

	return pm.init()
}

func GetPState(stats States, pluginName string) (*State, error) {
	var pstat *State

	for i, stat := range stats {
		if stat.Name == pluginName {
			pstat = &stats[i]
			break
		}
	}

	if pstat == nil {
		return nil, fmt.Errorf("%s does not exist", pluginName)
	}
	return pstat, nil
}

func (p Manager) configFilePath() string {
	configFilePath := path.Join(p.Storage, p.ConfigFile)
	return configFilePath
}

func (p Manager) pluginStorePath(pluginName string) string {
	return path.Join(p.Storage, pluginName)
}

func (p Manager) pluginJobPath(pluginName string) string {
	return path.Join(p.Job, pluginName)
}

func (p Manager) pluginBPMPath(pluginName string) string {
	return path.Join(p.Job, pluginName, "config", "bpm.yml")
}

func (p Manager) pluginBPMPathInStore(pluginName string) string {
	return path.Join(p.Storage, pluginName, "config", "bpm.yml")
}

func (p Manager) pluginMonitPath(pluginName string) string {
	return path.Join(p.Monit, fmt.Sprintf("%s.monitrc", pluginName))
}

func (p Manager) init() (Manager, error) {
	if _, err := os.Stat(p.configFilePath()); os.IsNotExist(err) {
		log.Println("initializing " + p.ConfigFile)
	}

	fh, err := os.OpenFile(p.configFilePath(), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return p, err
	}
	fh.Close()

	return p, p.bootstrapPlugin()
}

func (p Manager) bootstrapPlugin() error {
	stats, err := p.ListPlugins()
	if err != nil {
		return err
	}

	// recreate monitrc files for plugin. In case of deployment update, the custom monitrc files inserted by
	// plugin manager will be cleared.
	for _, stat := range stats {
		if err := monit.CreateMonitrcFor(stat.Name, p.pluginMonitPath(stat.Name)); err != nil {
			return err
		}
	}

	// recreate jobs folder for plugin. In case of stemcell update, the custom job folder will be destroyed.
	for _, stat := range stats {
		if _, err := os.Stat(p.pluginJobPath(stat.Name)); os.IsNotExist(err) {
			if err := p.copyPluginFromStorageToJob(stat.Name); err != nil {
				return err
			}
		}
	}

	// reload monit to pick up any monitrc updates
	if err := monit.Reload(); err != nil {
		return err
	}

	// start all enabled plugins
	// TODO: parallelize for efficiency
	for _, stat := range stats {
		if err := monit.Start(stat.Name); err != nil {
			return err
		}
	}

	return nil
}

type Application struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Command     string            `yaml:"command"`
	Memory      string            `yaml:"memory"`
	Args        []string          `yaml:"args"`
	Env         map[string]string `yaml:"env"`
}

type Info struct {
	Applications []Application `yaml:"applications"`
}

func getPluginInfo(s string) (Info, error) {
	log.Println("getting plugin info")
	info := Info{}

	r, err := zip.OpenReader(s)
	if err != nil {
		return info, errors.Wrap(err, "failed to open plugin file")
	}
	defer r.Close()

	for _, f := range r.File {
		log.Println(f.Name)
		if f.Name != PLUGIN_INFO_FILE {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return info, errors.Wrap(err, "failed to read plugin info")
		}
		err = yaml.NewDecoder(rc).Decode(&info)
		if err != nil {
			return info, errors.Wrap(err, "failed to unmarshal plugin info")
		}

		rc.Close()

		// TODO: verify only one application entry is provided
		return info, ValidatePluginName(info.Applications[0].Name)
	}
	return info, fmt.Errorf(PLUGIN_INFO_FILE + " does not exist in the plugin file")
}

var IsValidName = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`).MatchString

func ValidatePluginName(name string) error {
	if !IsValidName(name) {
		return fmt.Errorf("%s is not a valid plugin name, only 'a-z A-Z 0-9 _ -' are allowed", name)
	}
	return nil
}

func WriteYamlStructToFile(src interface{}, filename string) error {
	outContent, err := yaml.Marshal(src)
	if err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(outContent)
	return err
}

func ReadYamlStructFromFile(filename string, dst interface{}) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, dst)
}
